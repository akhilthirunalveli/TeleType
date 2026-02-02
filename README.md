# TeleType - Retro Terminal Chat System
![License](https://img.shields.io/badge/license-MIT-green.svg) ![Go](https://img.shields.io/badge/go-1.24-blue.svg)

> A retro-style terminal chat application built using **WebSockets** that supports real-time multi-client messaging, concurrent connections, and a classic interface inspired by early BBS systems.

## Features
- **Retro Theme**: High-fidelity CRT aesthetics, scanlines, and phosphor glow.
- **Web Client**: Zero-install web access with immersive CSS design.
- **Real-time**: Instant messaging powered by WebSockets.
- **Rooms & Identity**: 
  - Dynamic room switching (`/join <room>`).
  - Stateful usernames (`/setname <name>`).
- **Slash Commands**: Autocomplete menu for power users.
- **High Concurrency**: Supports thousands of concurrent users via Go goroutines.
- **Deployment Ready**: Includes Docker support and deployment guides.

## Prerequisites
- Go 1.21+ (for local build)
- Docker (optional, for containerized run)

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
