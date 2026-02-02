# Deploying TeleType

Since TeleType is a real-time WebSocket application that holds state in memory (users in rooms), it is best hosted on a **Virtual Private Server (VPS)** or a container platform that supports long-lived TCP/WebSocket connections.

**Note:** "Serverless" functions (like Vercel/AWS Lambda) are **NOT** recommended for this specific implementation because they kill connections frequently and don't share memory between instances.

## Option 1: Docker on a VPS (Recommended)
**Providers:** DigitalOcean ($4/mo), Hetzner, Linode, AWS EC2.

1.  **Get a Server**: Ubuntu Linux is standard.
2.  **SSH into Server**: `ssh root@<your-server-ip>`
3.  **Install Docker**:
    ```bash
    curl -fsSL https://get.docker.com | sh
    ```
4.  **Upload Code**:
    *   **Easy Way**: Clone your git repo.
    *   **Manual Way**: SCP your files.
    ```bash
    git clone https://github.com/yourusername/TeleType.git
    cd TeleType
    ```
5.  **Run with Compose**:
    ```bash
    docker compose up -d --build
    ```
6.  **Access**: `http://<your-server-ip>:8080`

## Option 2: Railway / Render / Fly.io (PaaS)
These platforms build your Dockerfile automatically.

### Railway / Render
1.  Connect your GitHub repository.
2.  It will detect the `Dockerfile`.
3.  **Important**: Set the internal port to `8080`.
4.  They will give you a domain like `teletype-production.up.railway.app`.

**Config Note:** Ensure the platform supports "WebSockets" (most do by default now).

## Option 3: Quick Sharing / Demo (Ngrok)
If you want to show a friend *right now* without buying a server.

1.  Start TeleType locally (`./server.exe`).
2.  Install [Ngrok](https://ngrok.com/).
3.  Run: `ngrok http 8080`
4.  It will give you a public URL (e.g., `https://a1b2.ngrok.io`).
5.  Send that URL to your friend.

## Important: Production Tweaks
If deploying to a real domain (e.g., `chat.example.com`):
1.  **HTTPS/WSS**:
    *   If using **Docker/VPS**, use Caddy or Nginx as a reverse proxy to handle SSL.
    *   **Caddy Example**:
        ```text
        chat.example.com {
            reverse_proxy localhost:8080
        }
        ```
    *   If using **Railway/Render**, they handle HTTPS automatically.
2.  **App Config**:
    *   Update `web/app.js` logic is already smart (`window.location.host`), so it should work automatically on any domain!
