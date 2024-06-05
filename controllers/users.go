package controllers

import (
	"eskimoe-server/database"
	"time"

	"github.com/gofiber/fiber/v2"
)

func JoinServer(c *fiber.Ctx) error {
	db := database.Database

	// Takes in: UniqueID, UniqueToken, DisplayName
	// Returns: Member
	member := new(struct {
		UniqueID    string `json:"unique_id"`
		UniqueToken string `json:"unique_token"`
		DisplayName string `json:"display_name"`
	})

	if err := c.BodyParser(member); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errorCode": fiber.StatusBadRequest,
			"error":     "Invalid Request Body",
		})
	}

	// Empty Values Check
	if member.UniqueID == "" || member.UniqueToken == "" || member.DisplayName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errorCode": fiber.StatusBadRequest,
			"error":     "One or more fields are empty",
		})
	}

	// Check if Member Exists
	var existingMember database.Member
	if db.Where("unique_id = ?", member.UniqueID).First(&existingMember).Error == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errorCode": fiber.StatusBadRequest,
			"error":     "Member already exists",
		})
	}

	// Create Member
	joinedAt := time.Now().Format(time.RFC3339)
	newMember := database.Member{
		UniqueID:    member.UniqueID,
		UniqueToken: member.UniqueToken,
		DisplayName: member.DisplayName,
		JoinedAt:    joinedAt,
	}

	if err := db.Create(&newMember).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errorCode": fiber.StatusInternalServerError,
			"error":     "Error Creating Member",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Member Created",
		"member":  newMember,
	})
}

func GetMember(c *fiber.Ctx) error {
	db := database.Database

	// Takes in: UniqueID
	// Returns: Member
	member := new(struct {
		UniqueID string `json:"unique_id"`
	})

	if err := c.BodyParser(member); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errorCode": fiber.StatusBadRequest,
			"error":     "Invalid Request Body",
		})
	}

	// Empty Values Check
	if member.UniqueID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errorCode": fiber.StatusBadRequest,
			"error":     "One or more fields are empty",
		})
	}

	// Check if Member Exists
	var existingMember database.Member
	if db.Where("unique_id = ?", member.UniqueID).First(&existingMember).Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errorCode": fiber.StatusBadRequest,
			"error":     "Member does not exist",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"member": existingMember,
	})
}
