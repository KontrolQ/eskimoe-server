package middleware

// Auth is simple. If a token is provided in the Authorization header, it will be checked against the database.
// If the token is valid, the member is attached to the context as a local variable.
// If the token is invalid, the request is still forwarded, but the member is nil.
// The following function will decide what to do with the member.

import (
	"eskimoe-server/database"

	"github.com/gofiber/fiber/v2"
)

func Auth(c *fiber.Ctx) error {
	// Get the Authorization header
	token := c.Get("Authorization")

	// If the token is empty, continue
	if token == "" {
		return c.Next()
	}

	// Check the token against the database
	db := database.Database

	var member database.Member

	if db.Where("auth_token = ?", token).Preload("Roles").First(&member).Error != nil {
		return c.Next()
	}

	c.Locals("Member", member)

	return c.Next()
}
