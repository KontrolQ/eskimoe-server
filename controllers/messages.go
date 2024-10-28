package controllers

import (
	"encoding/json"
	"eskimoe-server/config"
	"eskimoe-server/database"
	"eskimoe-server/socket"

	"github.com/gofiber/fiber/v2"
)

// Gets last 25 messages from the room passed in the URL
func GetMessages(c *fiber.Ctx) error {
	_, err := c.Locals("Member").(database.Member)

	if !err {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"errorCode": fiber.StatusUnauthorized,
			"error":     "Unauthorized",
		})
	}

	db := database.Database

	roomID := c.Params("room")

	var room database.Room

	if err := db.Preload("Messages").Preload("Messages.Author").Preload("Messages.Reactions").Preload("Messages.Attachments").Where("id = ?", roomID).Order("created_at desc").Limit(25).Find(&room).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"errorCode": fiber.StatusNotFound,
			"error":     "Room Not Found",
		})
	} else {
		return c.Status(fiber.StatusOK).JSON(room.Messages)
	}
}

// Send a message to the room passed in the URL
func SendMessage(c *fiber.Ctx) error {
	member, ok := c.Locals("Member").(database.Member)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"errorCode": fiber.StatusUnauthorized,
			"error":     "Unauthorized",
		})
	}

	db := database.Database

	roomID := c.Params("room")

	var room database.Room
	if err := db.First(&room, roomID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"errorCode": fiber.StatusNotFound,
			"error":     "Room Not Found",
		})
	}

	var message database.Message
	if err := c.BodyParser(&message); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errorCode": fiber.StatusBadRequest,
			"error":     "Bad Request",
		})
	}

	db.Model(&message).Association("Author").Append(&member)
	db.Model(&message).Association("Room").Append(&room)

	if err := db.Create(&message).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errorCode": fiber.StatusInternalServerError,
			"error":     "Error Creating Message",
		})
	}

	broadcastData, err := json.Marshal(config.SocketBroadcast{
		BroadcastType: config.MessageCreated,
		Data:          message,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errorCode": fiber.StatusInternalServerError,
			"error":     "Error Encoding Message",
		})
	}

	socket.WsHub.Broadcast <- broadcastData

	return c.Status(fiber.StatusCreated).JSON(message)
}

func DeleteMessage(c *fiber.Ctx) error {
	deleter, err := c.Locals("Member").(database.Member)

	if !err {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"errorCode": fiber.StatusUnauthorized,
			"error":     "Unauthorized",
		})
	}

	db := database.Database

	messageID := c.Params("message")
	roomID := c.Params("room")

	var message database.Message

	if err := db.Preload("Author").Where("id = ? AND room_id = ?", messageID, roomID).First(&message).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"errorCode": fiber.StatusNotFound,
			"error":     "Message Not Found",
		})
	}

	// Must be the author, owner, or admin to delete a message
	hasDeletePermission := message.Author.ID == deleter.ID

	if !hasDeletePermission {
		if config.Owner == deleter.UniqueID {
			hasDeletePermission = true
		} else {
			for _, role := range deleter.Roles {
				permissions := role.Permissions
				for _, permission := range permissions {
					if permission == database.DeleteMessage {
						hasDeletePermission = true
						break
					}
				}
			}
		}
	}

	if !hasDeletePermission {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"errorCode": fiber.StatusUnauthorized,
			"error":     "Unauthorized",
		})
	}

	if err := db.Delete(&message).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errorCode": fiber.StatusInternalServerError,
			"error":     "Error Deleting Message",
		})
	}

	deletedData := struct {
		MessageID int  `json:"message_id"`
		RoomID    int  `json:"room_id"`
		Deleted   bool `json:"deleted"`
	}{
		MessageID: message.ID,
		RoomID:    message.RoomID,
		Deleted:   true,
	}

	broadcastData, _ := json.Marshal(config.SocketBroadcast{
		BroadcastType: config.MessageDeleted,
		Data:          deletedData,
	})

	socket.WsHub.Broadcast <- broadcastData

	return c.Status(fiber.StatusOK).JSON(deletedData)
}
