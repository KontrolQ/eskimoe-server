package controllers

import (
	"eskimoe-server/database"
	"eskimoe-server/utils"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func CategoryWiseRooms(c *fiber.Ctx) error {
	_, err := c.Locals("Member").(database.Member)

	if !err {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"errorCode": fiber.StatusUnauthorized,
			"error":     "Unauthorized",
		})
	}

	db := database.Database

	var categories []database.Category

	if err := db.Preload("Rooms").Find(&categories).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errorCode": fiber.StatusInternalServerError,
			"error":     "Error Finding Categories",
		})
	}

	return c.Status(fiber.StatusOK).JSON(categories)
}

func CreateRoom(c *fiber.Ctx) error {
	member, err := c.Locals("Member").(database.Member)

	if !err {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"errorCode": fiber.StatusUnauthorized,
			"error":     "Unauthorized",
		})
	}
	db := database.Database

	if !utils.VerifyOwnerOrPermission(member, "manage_rooms") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"errorCode": fiber.StatusUnauthorized,
			"error":     "Unauthorized",
		})
	}

	roomCreationStruct := new(struct {
		Name        string            `json:"name"`
		Description string            `json:"description"`
		CategoryID  int               `json:"category_id"`
		Type        database.RoomType `json:"type"`
	})

	if err := c.BodyParser(roomCreationStruct); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errorCode": fiber.StatusBadRequest,
			"error":     "Invalid Request",
		})
	}

	newRoom := database.Room{
		Name:        roomCreationStruct.Name,
		Description: roomCreationStruct.Description,
		CategoryID:  roomCreationStruct.CategoryID,
		Type:        roomCreationStruct.Type,
	}

	if err := db.Create(&newRoom).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errorCode": fiber.StatusInternalServerError,
			"error":     "Error Creating Room",
		})
	}

	// Update Category Room Order
	var category database.Category

	if err := db.First(&category, roomCreationStruct.CategoryID).Error; err != nil {
		if err := db.Delete(&newRoom).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"errorCode": fiber.StatusInternalServerError,
				"error":     "Error Deleting Room",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errorCode": fiber.StatusInternalServerError,
			"error":     "Error Finding Category",
		})
	}

	category.RoomOrder = append(category.RoomOrder, newRoom.ID)

	if err := db.Save(&category).Error; err != nil {
		if err := db.Delete(&newRoom).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"errorCode": fiber.StatusInternalServerError,
				"error":     "Error Deleting Room",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errorCode": fiber.StatusInternalServerError,
			"error":     "Error Updating Category",
		})
	}

	// Update the Server Log
	serverLog := database.Log{
		Type:     database.RoomCreated,
		Content:  fmt.Sprintf("Room %s created in Category %s", newRoom.Name, category.Name),
		MemberID: member.ID,
		ServerID: member.ServerID,
	}

	if err := db.Create(&serverLog).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errorCode": fiber.StatusInternalServerError,
			"error":     "Error Creating Log",
		})
	}

	return c.Status(fiber.StatusOK).JSON(newRoom)
}

func UpdateRoom(c *fiber.Ctx) error {
	member, err := c.Locals("Member").(database.Member)

	if !err {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"errorCode": fiber.StatusUnauthorized,
			"error":     "Unauthorized",
		})
	}

	db := database.Database

	if !utils.VerifyOwnerOrPermission(member, "manage_rooms") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"errorCode": fiber.StatusUnauthorized,
			"error":     "Unauthorized",
		})
	}

	roomUpdateStruct := new(struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	})

	roomID := c.Params("room")

	if err := c.BodyParser(roomUpdateStruct); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errorCode": fiber.StatusBadRequest,
			"error":     "Invalid Request",
		})
	}

	var room database.Room
	var changes []string

	if err := db.First(&room, roomID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"errorCode": fiber.StatusNotFound,
			"error":     "Room Not Found",
		})
	}

	if roomUpdateStruct.Name != "" && room.Name != roomUpdateStruct.Name {
		room.Name = roomUpdateStruct.Name
		changes = append(changes, fmt.Sprintf("Name: %s", room.Name))
	}

	if roomUpdateStruct.Description != "" && room.Description != roomUpdateStruct.Description {
		room.Description = roomUpdateStruct.Description
		changes = append(changes, fmt.Sprintf("Description: %s", room.Description))
	}

	if err := db.Save(&room).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errorCode": fiber.StatusInternalServerError,
			"error":     "Error Updating Room",
		})
	}

	// Update the Server Log on changes
	if len(changes) == 0 {
		return c.Status(fiber.StatusOK).JSON(room)
	}

	serverLog := database.Log{
		Type:     database.RoomUpdated,
		Content:  fmt.Sprintf("Room %s updated.\n%s", room.Name, strings.Join(changes, "\n")),
		MemberID: member.ID,
		ServerID: member.ServerID,
	}

	if err := db.Create(&serverLog).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errorCode": fiber.StatusInternalServerError,
			"error":     "Error Creating Log",
		})
	}

	return c.Status(fiber.StatusOK).JSON(room)
}

func DeleteRoom(c *fiber.Ctx) error {
	member, err := c.Locals("Member").(database.Member)

	if !err {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"errorCode": fiber.StatusUnauthorized,
			"error":     "Unauthorized",
		})
	}

	db := database.Database

	if !utils.VerifyOwnerOrPermission(member, "manage_rooms") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"errorCode": fiber.StatusUnauthorized,
			"error":     "Unauthorized",
		})
	}

	roomID := c.Params("room")

	var room database.Room

	if err := db.First(&room, roomID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"errorCode": fiber.StatusNotFound,
			"error":     "Room Not Found",
		})
	}

	if err := db.Delete(&room).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errorCode": fiber.StatusInternalServerError,
			"error":     "Error Deleting Room",
		})
	}

	// Update Category Room Order
	var category database.Category

	if err := db.First(&category, room.CategoryID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errorCode": fiber.StatusInternalServerError,
			"error":     "Error Finding Category",
		})
	}

	for i, roomID := range category.RoomOrder {
		if roomID == room.ID {
			category.RoomOrder = append(category.RoomOrder[:i], category.RoomOrder[i+1:]...)
			break
		}
	}

	if err := db.Save(&category).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errorCode": fiber.StatusInternalServerError,
			"error":     "Error Updating Category",
		})
	}

	// Update the Server Log
	serverLog := database.Log{
		Type:     database.RoomDeleted,
		Content:  fmt.Sprintf("Room %s deleted from Category %s", room.Name, category.Name),
		MemberID: member.ID,
		ServerID: member.ServerID,
	}

	if err := db.Create(&serverLog).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errorCode": fiber.StatusInternalServerError,
			"error":     "Error Creating Log",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Room Deleted",
	})
}
