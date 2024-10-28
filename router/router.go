package router

import (
	"eskimoe-server/controllers"
	"eskimoe-server/database"
	"eskimoe-server/socket"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func Initialize(router *fiber.App) {
	// Server Endpoints
	router.Get("/", controllers.ServerInfo)

	members := router.Group("/members")

	// Members Endpoints
	members.Post("/join", controllers.JoinServer)
	members.Delete("/leave", controllers.LeaveServer)
	members.Get("/me", controllers.Me)
	members.Post("/me", controllers.Me)

	// Rooms Endpoints
	rooms := router.Group("/rooms")

	rooms.Get("/", controllers.CategoryWiseRooms)
	rooms.Post("/new", controllers.CreateRoom)
	rooms.Patch("/:room", controllers.UpdateRoom)
	rooms.Delete("/:room", controllers.DeleteRoom)

	// Messages Endpoints
	messages := rooms.Group("/:room/messages")

	messages.Get("/", controllers.GetMessages)
	messages.Post("/new", controllers.SendMessage)
	messages.Delete("/:message", controllers.DeleteMessage)

	router.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("SocketCapable", true)
			return c.Next()
		}
		return c.Status(fiber.StatusUpgradeRequired).JSON(fiber.Map{
			"errorCode": fiber.StatusUpgradeRequired,
			"error":     "Socket Upgrade Required.",
		})
	})

	router.Get("/ws/listen", websocket.New(func(c *websocket.Conn) {
		if _, ok := c.Locals("Member").(database.Member); !ok {
			log.Println("Unauthorized Member Disconnected")
			c.Close()
			return
		}

		log.Println("Connected Member", c.Locals("Member").(database.Member).DisplayName)

		socket.WsHub.Register <- c
		defer func() {
			socket.WsHub.Unregister <- c
		}()

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("Read Error:", err)
				return
			}

			// Handle Ping messages
			if string(msg) == string(rune(websocket.PingMessage)) {
				c.WriteMessage(websocket.PongMessage, []byte(string(rune(websocket.PongMessage))))
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
