package main

import (
	"github.com/coder/websocket"
	"log"
	"net/http"
	"teletype/internal/protocol"
	"teletype/internal/server"
)

func serveWs(hub *server.Hub, w http.ResponseWriter, r *http.Request) {
	// Accept the connection
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true, // Allow all origins for now (dev mode)
	})
	if err != nil {
		log.Printf("Error accepting websocket: %v", err)
		return
	}

	// Create a new client
	client := &server.Client{
		ID:   r.RemoteAddr, // Simple ID for now
		Hub:  hub,
		Conn: c,
		Send: make(chan protocol.Message, 256),
	}

	// Register with hub
	client.Hub.Register <- client

	// Start pumps
	// Note: In coder/websocket, typically we just need a read loop that blocks.
	// Write loop is separate.
	// We need to manage contexts properly.

	// We'll spin off write pump in a goroutine
	go client.WritePump(r.Context())

	// Read pump blocks until connection closes
	client.ReadPump(r.Context())
}

func main() {
	log.Println("Starting TeleType Server...")

	hub := server.NewHub()
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	// Serve static files
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	log.Println("Server listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
