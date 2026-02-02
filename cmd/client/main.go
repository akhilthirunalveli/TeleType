package main

import (
	"context"
	"flag"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"log"
	"os"
	"teletype/internal/protocol"
	"teletype/internal/ui"
	"time"
)

func main() {
	serverAddr := flag.String("addr", "ws://localhost:8080/ws", "server address")
	username := flag.String("user", "Guest", "username")
	room := flag.String("room", "general", "room to join")
	flag.Parse()

	termUI := ui.NewUI()
	termUI.Init()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// Connect to server
	c, _, err := websocket.Dial(ctx, *serverAddr, nil)
	if err != nil {
		log.Fatal("Dial:", err)
	}
	defer c.Close(websocket.StatusNormalClosure, "")

	// Reset context for long running connection
	ctx = context.Background()

	// Join logic
	joinMsg := protocol.Message{
		Type:    protocol.MsgTypeJoin,
		Content: *room, // room name
		Sender:  *username,
	}
	wsjson.Write(ctx, c, joinMsg)

	// Start read loop
	go func() {
		for {
			var msg protocol.Message
			err := wsjson.Read(ctx, c, &msg)
			if err != nil {
				log.Printf("Read error: %v", err)
				os.Exit(1)
				return
			}
			termUI.PrintMessage(msg)
		}
	}()

	// Input loop
	for {
		text := termUI.Prompt()
		if text == "" {
			continue // Handle EOF or empty
		}
		if text == "/quit" {
			break
		}

		msg := protocol.NewChatMessage(*username, text, *room)
		// We set Sender client-side here for UI echo purposes maybe,
		// but server enforces it usually.
		// Ideally server shouldn't trust client sender ID, but for this demo
		// we pass it in join or just assume server tracks it by connection.
		// The protocol definition has Sender.

		err := wsjson.Write(ctx, c, msg)
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}
