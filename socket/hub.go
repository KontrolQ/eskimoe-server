package socket

import (
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type Hub struct {
	Clients    map[*websocket.Conn]bool
	Broadcast  chan []byte
	Register   chan *websocket.Conn
	Unregister chan *websocket.Conn
	mu         sync.Mutex
}

var WsHub = Hub{
	Clients:    make(map[*websocket.Conn]bool),
	Broadcast:  make(chan []byte),
	Register:   make(chan *websocket.Conn),
	Unregister: make(chan *websocket.Conn),
}

func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.Register:
			h.mu.Lock()
			h.Clients[conn] = true
			h.mu.Unlock()
		case conn := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[conn]; ok {
				delete(h.Clients, conn)
				conn.Close()
			}
			h.mu.Unlock()
		case message := <-h.Broadcast:
			h.mu.Lock()
			for conn := range h.Clients {
				if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
					conn.Close()
					delete(h.Clients, conn)
				}
			}
			h.mu.Unlock()
		}
	}
}
