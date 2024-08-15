package main

import (
	"log"

	"eskimoe-server/config"
	"eskimoe-server/database"
	"eskimoe-server/middleware"
	"eskimoe-server/router"
	"eskimoe-server/socket"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	database.Initialize()

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowHeaders: "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin,Authorization",
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))
	app.Use(logger.New())

	app.Use(middleware.Json)
	app.Use(middleware.Auth)

	app.Use(helmet.New())

	go socket.WsHub.Run()

	router.Initialize(app)

	log.Fatal(app.Listen(":" + config.Port))
}
