package main

import (
	"log"

	"eskimoe-server/config"
	"eskimoe-server/database"
	"eskimoe-server/middleware"
	"eskimoe-server/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowHeaders: "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin,Authorization",
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))
	app.Use(logger.New())

	app.Use(middleware.Json)

	app.Use(helmet.New())

	router.Initialize(app)

	database.Initialize()

	log.Fatal(app.Listen(":" + config.Port))
}
