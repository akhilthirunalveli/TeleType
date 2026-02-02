package server

import (
	"log"
	"sync"
	"teletype/internal/protocol"
)

// Hub maintains the set of active clients and broadcasts messages to the
// rooms.
type Hub struct {
	// Registered clients.
	Clients map[*Client]bool

	// Rooms maps room names to a set of clients
	Rooms map[string]map[*Client]bool

	// Inbound messages from the clients.
	Broadcast chan protocol.Message

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client

	mu sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan protocol.Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Rooms:      make(map[string]map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client] = true
			h.mu.Unlock()
			log.Printf("Client connected: %s", client.ID)

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				client.Close()
				// Remove from rooms
				if client.Room != "" {
					if clients, ok := h.Rooms[client.Room]; ok {
						delete(clients, client)
						if len(clients) == 0 {
							delete(h.Rooms, client.Room)
						}
					}
				}
			}
			h.mu.Unlock()
			log.Printf("Client disconnected: %s", client.ID)

		case message := <-h.Broadcast:
			h.mu.RLock()
			// If message has a room, broadcast only to that room
			if message.Room != "" {
				if clients, ok := h.Rooms[message.Room]; ok {
					for client := range clients {
						select {
						case client.Send <- message:
						default:
							close(client.Send)
							delete(h.Clients, client)
						}
					}
				}
			} else {
				// Global broadcast (if ever needed)
				for client := range h.Clients {
					select {
					case client.Send <- message:
					default:
						close(client.Send)
						delete(h.Clients, client)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

// JoinRoom handles adding a client to a specific room
func (h *Hub) JoinRoom(client *Client, roomName string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Leave current room if any
	if client.Room != "" {
		if clients, ok := h.Rooms[client.Room]; ok {
			delete(clients, client)
			if len(clients) == 0 {
				delete(h.Rooms, client.Room)
			}
		}
	}

	client.Room = roomName
	if _, ok := h.Rooms[roomName]; !ok {
		h.Rooms[roomName] = make(map[*Client]bool)
	}
	h.Rooms[roomName][client] = true
}
