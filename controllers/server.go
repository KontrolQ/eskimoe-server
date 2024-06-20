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

	// Get Server ID 1
	var server database.Server

	if err := db.First(&server).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(defaultResponse)
	}

	categories := []database.Category{}

	if err := db.Find(&categories).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(defaultResponse)
	}

	for i, category := range categories {
		room := []database.Room{}

		if err := db.Where("category_id = ?", category.ID).Find(&room).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(defaultResponse)
		}

		categories[i].Rooms = append(categories[i].Rooms, room...)
	}

	reactions := []database.ServerReaction{}

	if err := db.Find(&reactions).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(defaultResponse)
	}

	roles := []database.Role{}

	if err := db.Find(&roles).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(defaultResponse)
	}

	members := []database.Member{}

	if err := db.Find(&members).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(defaultResponse)
	}

	// attach roles, remove auth token
	for i, member := range members {
		role := []database.Role{}

		if err := db.Model(&member).Association("Roles").Find(&role); err != nil {
			return c.Status(fiber.StatusNotFound).JSON(defaultResponse)
		}

		members[i].Roles = append(members[i].Roles, role...)
		members[i].AuthToken = ""
	}

	server.Categories = append(server.Categories, categories...)
	server.ServerReactions = append(server.ServerReactions, reactions...)
	server.Roles = append(server.Roles, roles...)
	server.Members = append(server.Members, members...)

	return c.Status(fiber.StatusOK).JSON(server)
}
