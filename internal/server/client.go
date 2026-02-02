package server

import (
	"context"
	"log"
	"teletype/internal/protocol"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	ID   string
	Name string
	Hub  *Hub
	Conn *websocket.Conn
	Send chan protocol.Message
	Room string
}

// ReadPump pumps messages from the websocket connection to the hub.
func (c *Client) ReadPump(ctx context.Context) {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close(websocket.StatusNormalClosure, "")
	}()

	c.Conn.SetReadLimit(maxMessageSize)

	for {
		var msg protocol.Message
		err := wsjson.Read(ctx, c.Conn, &msg)
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
				websocket.CloseStatus(err) == websocket.StatusGoingAway {
				break
			}
			log.Printf("error reading json: %v", err)
			break
		}

		// Enforce sender ID consistency
		// If Name is set, use it; otherwise fall back to ID (IP)
		if c.Name != "" {
			msg.Sender = c.Name
		} else {
			msg.Sender = c.ID
		}
		msg.Timestamp = time.Now()

		// Handle special messages (Join room, etc) or just forward chat
		if msg.Type == protocol.MsgTypeJoin {
			c.Hub.JoinRoom(c, msg.Content)
			// Notify room
			sysMsg := protocol.NewSystemMessage(c.ID + " joined room " + msg.Content)
			if c.Name != "" {
				sysMsg.Content = c.Name + " joined room " + msg.Content
			}
			sysMsg.Room = msg.Content
			c.Hub.Broadcast <- sysMsg
			msg.Room = msg.Content
		}

		if msg.Type == protocol.MsgTypeName {
			oldName := c.Name
			if oldName == "" {
				oldName = c.ID
			}
			c.Name = msg.Content
			// Notify room
			if c.Room != "" {
				sysMsg := protocol.NewSystemMessage(oldName + " changed name to " + c.Name)
				sysMsg.Room = c.Room
				c.Hub.Broadcast <- sysMsg
			}
			continue // Don't broadcast the NAME message itself
		}

		if msg.Room == "" && c.Room != "" {
			msg.Room = c.Room
		}

		// Force room to be client's current room to avoid spoofing or leakage
		if c.Room != "" {
			msg.Room = c.Room
		}

		// Only broadcast if it's a chat message
		if msg.Type == protocol.MsgTypeChat {
			c.Hub.Broadcast <- msg
		}
	}
}

// WritePump pumps messages from the hub to the websocket connection.
func (c *Client) WritePump(ctx context.Context) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close(websocket.StatusNormalClosure, "")
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				// The hub closed the channel.
				c.Conn.Close(websocket.StatusNormalClosure, "Channel closed")
				return
			}

			// Add a deadline for writing
			ctx, cancel := context.WithTimeout(ctx, writeWait)
			err := wsjson.Write(ctx, c.Conn, message)
			cancel()
			if err != nil {
				return
			}

		case <-ticker.C:
			ctx, cancel := context.WithTimeout(ctx, writeWait)
			err := c.Conn.Ping(ctx)
			cancel()
			if err != nil {
				return
			}

		case <-ctx.Done():
			return
		}
	}
}

func (c *Client) Close() {
	close(c.Send)
}
