package controllers

import (
	"eskimoe-server/config"
	"eskimoe-server/database"

	"github.com/gofiber/fiber/v2"
)

func ServerInfo(c *fiber.Ctx) error {
	_, err := c.Locals("Member").(database.Member)

	defaultResponse := fiber.Map{
		"name":    config.Name,
		"message": config.Message,
		"version": config.Version,
	}

	if !err {
		return c.Status(fiber.StatusOK).JSON(defaultResponse)
	}

	db := database.Database

	var server database.Server

	if err := db.Model(&database.Server{}).
		Preload("Categories").
		Preload("Categories.Rooms").
		Preload("ServerReactions").
		Preload("Roles").
		Preload("Events").
		Preload("Members.Roles").
		Preload("Members").
		First(&server).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(defaultResponse)
	}

	// Remove the auth token from all members
	for i := range server.Members {
		server.Members[i].AuthToken = ""
	}

	return c.Status(fiber.StatusOK).JSON(server)
}
