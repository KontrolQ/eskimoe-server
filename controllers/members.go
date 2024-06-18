package controllers

import (
	"eskimoe-server/database"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func JoinServer(c *fiber.Ctx) error {
	db := database.Database

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

	// encrypt the token
	encryptedToken, err := bcrypt.GenerateFromPassword([]byte(member.UniqueToken), bcrypt.DefaultCost)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errorCode": fiber.StatusInternalServerError,
			"error":     "Error Encrypting Token",
		})
	}

	// Create Member
	joinedAt := time.Now()
	newMember := database.Member{
		UniqueID:    member.UniqueID,
		UniqueToken: member.UniqueToken,
		AuthToken:   string(encryptedToken),
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

func LeaveServer(c *fiber.Ctx) error {
	// Get the member from the context
	member, err := c.Locals("Member").(database.Member)

	// UnAuthorized
	if !err {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"errorCode": fiber.StatusUnauthorized,
			"error":     "Unauthorized",
		})
	}

	db := database.Database

	// Delete the member
	if err := db.Delete(&member).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errorCode": fiber.StatusInternalServerError,
			"error":     "Error Leaving Server",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Left Server",
	})
}

// Changes the about information of the member and returns the updated member
func Me(c *fiber.Ctx) error {
	member, err := c.Locals("Member").(database.Member)

	// UnAuthorized
	if !err {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"errorCode": fiber.StatusUnauthorized,
			"error":     "Unauthorized",
		})
	}

	db := database.Database

	newMember := database.Member{}

	if err := c.BodyParser(&newMember); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errorCode": fiber.StatusBadRequest,
			"error":     "Invalid Request Body",
		})
	}

	// Update the member whatever is provided
	if newMember.About != "" {
		member.About = newMember.About
	}

	if newMember.Pronouns != "" {
		member.Pronouns = newMember.Pronouns
	}

	if newMember.DisplayName != "" {
		member.DisplayName = newMember.DisplayName
	}

	if err := db.Save(&member).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errorCode": fiber.StatusInternalServerError,
			"error":     "Error Updating Member",
		})
	}

	// Remove the auth_token field from the JSON
	member.AuthToken = ""

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Member Updated",
		"member":  member,
	})
}
