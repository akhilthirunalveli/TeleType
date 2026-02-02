# TeleType - Retro Terminal Chat System

A high-concurrency, retro-styled terminal chat system built with Go and WebSockets.

## Features
- **Retro Theme**: 90s green-screen terminal UI with ANSI graphics.
- **Real-time**: WebSocket-based communication.
- **Concurrency**: Handles thousands of concurrent connections (goroutine-per-client).
- **Rooms**: Join different channels/rooms.
- **Dockerized**: Easy deployment with Docker Compose.

## Prerequisites
- Go 1.24+
- Docker (optional)

## Building Locally

```bash
# Build server
go build -o server ./cmd/server

# Build client
go build -o client ./cmd/client
```

## Running

### 1. Start Server
You can run the server directly or via Docker.

**Directly:**
```bash
./server
```
Server listens on port `:8080`.

**Docker:**
```bash
docker-compose up --build
```

### 2. Connect Client
Open multiple terminal windows to simulate users.

```bash
# Default (Guest, general room)
./client

# Custom user and room
./client -user Alice -room devops
./client -user Bob -room devops
```

## LAN / Offline Usage

To use TeleType without internet on a local network:

1.  **Find Server IP**: Run `ipconfig` (Windows) or `ifconfig` (Linux/Mac) on the server machine to get its LAN IP (e.g., `192.168.1.5`).
2.  **Start Server**: `./server.exe` (Server listens on all interfaces by default).
3.  **Connect Client**:
    ```bash
    ./client.exe -addr ws://192.168.137.1:8080/ws -user Guest -room general
    ```

To build without internet, we have vendored dependencies. Use `go build -mod=vendor` if needed (though Go defaults to vendor if present in recent versions).

## Protocol
TeleType uses JSON over WebSockets.
- `JOIN`: Join a room.
- `CHAT`: Standard message.
- `SYSTEM`: Server notifications.

## Architecture
- **Tech Stack**: Go (Standard Lib + coder/websocket), Docker.
- **Design**: Implements the Actor model (roughly) using goroutines for each client connection and a central Hub for routing.
