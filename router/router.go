package router

import (
	"eskimoe-server/config"
	"eskimoe-server/controllers"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func Initialize(router *fiber.App) {
	router.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"name":    config.Name,
			"message": config.Message,
			"version": config.Version,
		})
	})

	router.Post("/join", controllers.JoinServer)
	router.Get("/member", controllers.GetMember)

	// Not ready: Just a test
	router.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("SocketCapable", true)
			return c.Next()
		}
		return c.Status(fiber.StatusUpgradeRequired).JSON(fiber.Map{
			"errorCode": fiber.StatusUpgradeRequired,
			"error":     "Socket Upgrade Required.",
		})
	})

	router.Get("/ws/:room", websocket.New(func(c *websocket.Conn) {
		// c.Locals is added to the *websocket.Conn
		log.Println("SocketCapable:", c.Locals("SocketCapable"))
		log.Println("Room:", c.Params("room"))

		var (
			mt  int
			msg []byte
			err error
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Println("Error Reading Socket Message:", err)
				break
			}
			log.Printf("Message Received: %s", msg)

			response := "Echo: " + string(msg)
			if err = c.WriteMessage(mt, []byte(response)); err != nil {
				log.Println("Error Writing Socket Message:", err)
				break
			}
		}
	}))

	router.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"errorCode": fiber.StatusNotFound,
			"error":     "Unsupported Endpoint.",
		})
	})
}
