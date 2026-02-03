const display = document.getElementById('chat-display');
const input = document.getElementById('message-input');

let ws;
let username = "WebUser_" + Math.floor(Math.random() * 1000);
let room = "general";

function connect() {
    let wsUrl;
    if (window.location.protocol === 'file:') {
        wsUrl = "ws://localhost:8080/ws";
        addMessage("SYSTEM", "Running local file. Attempting localhost:8080...");
    } else {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        wsUrl = `${protocol}//${window.location.host}/ws`;
    }

    ws = new WebSocket(wsUrl);

    ws.onopen = () => {
        addMessage("SYSTEM", "Connected to mainframe.");
        // Send Join message
        ws.send(JSON.stringify({
            type: "JOIN",
            content: room,
            sender: username,
            room: room
        }));
    };

    ws.onerror = (e) => {
        console.error("WebSocket Error:", e);
        addMessage("SYSTEM", "Connection Error. Ensure server is running.", "ERROR");
    };

    ws.onmessage = (event) => {
        try {
            const msg = JSON.parse(event.data);
            addMessage(msg.sender, msg.content, msg.type, msg.timestamp);
        } catch (e) {
            console.error("Error parsing message:", e);
        }
    };

    ws.onclose = () => {
        addMessage("SYSTEM", "Signal lost. Retrying uplink...", "ERROR");
        setTimeout(connect, 3000);
    };
}

function addMessage(sender, content, type = "CHAT", timestamp = new Date()) {
    const line = document.createElement('div');
    line.className = 'msg-line';

    const tsStr = new Date(timestamp).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' });

    if (type === "SYSTEM" || type === "JOIN") {
        line.innerHTML = `<span class="system-msg">[${tsStr}] ${content}</span>`;
    } else {
        line.innerHTML = `<span class="msg-timestamp">[${tsStr}]</span> <span class="msg-sender">${sender}:</span> <span class="msg-content">${content}</span>`;
    }

    display.appendChild(line);
    display.scrollTop = display.scrollHeight;
}

// Parse URL params
const urlParams = new URLSearchParams(window.location.search);
if (urlParams.has('room')) room = urlParams.get('room');
if (urlParams.has('user')) username = urlParams.get('user');

function sendMessage() {
    const text = input.value.trim();
    if (!text) return;

    // Handle slash commands
    if (text.startsWith('/')) {
        const parts = text.split(' ');
        const cmd = parts[0].toLowerCase();

        if (cmd === '/join' && parts.length > 1) {
            const newRoom = parts[1];
            addMessage("SYSTEM", `Switching to room: ${newRoom}`);
            room = newRoom;
            ws.send(JSON.stringify({
                type: "JOIN",
                content: room,
                sender: username,
                room: room
            }));
            input.value = '';
            return;
        }

        if ((cmd === '/nick' || cmd === '/setname') && parts.length > 1) {
            username = parts[1];
            ws.send(JSON.stringify({
                type: "NAME",
                content: username,
                sender: username,
                room: room
            }));
            input.value = '';
            return;
        }
    }

    if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({
            type: "CHAT",
            content: text,
            sender: username,
            room: room,
        }));
        input.value = '';
        suggestionsBox.classList.add('hidden');
    }
}

input.addEventListener('keypress', (e) => {
    if (e.key === 'Enter') {
        sendMessage();
    }
});

// Keep focus on input unless selecting text
document.addEventListener('click', (e) => {
    if (window.getSelection().toString() === "") {
        input.focus();
    }
});

// Command Suggestions
const suggestionsBox = document.getElementById('command-suggestions');
const commands = [
    { cmd: '/join', desc: 'Join a room (Usage: /join <room>)' },
    { cmd: '/setname', desc: 'Set your name (Usage: /setname <name>)' },
    { cmd: '/nick', desc: 'Alias for /setname' },
    { cmd: '/help', desc: 'Show this help menu' } // We'll handle /help to show locally too if wanted
];

input.addEventListener('input', (e) => {
    const val = input.value;
    if (val.startsWith('/')) {
        suggestionsBox.classList.remove('hidden');
        // Filter based on input? Or just show all?
        // Let's filter slightly
        const search = val.toLowerCase();
        const matches = commands.filter(c => c.cmd.startsWith(search) || search === '/');

        if (matches.length > 0) {
            suggestionsBox.innerHTML = matches.map(c =>
                `<div class="suggestion-item"><span class="suggestion-cmd">${c.cmd}</span> ${c.desc}</div>`
            ).join('');
        }
    } else {
        suggestionsBox.classList.add('hidden');
    }
});

// Start
addMessage("SYSTEM", "Booting up TeleType Web Client...");
connect();

// Mobile Viewport Fix for Virtual Keyboard
function adjustViewport() {
    if (window.innerWidth <= 768) {
        // Set a CSS variable to the actual visible height
        document.documentElement.style.setProperty('--vh', `${window.innerHeight * 0.01}px`);
        
        // Ensure input is visible when keyboard opens
        if (document.activeElement === input) {
            setTimeout(() => {
                input.scrollIntoView({ behavior: 'smooth', block: 'center' });
            }, 300);
        }
    }
}

window.addEventListener('resize', adjustViewport);
window.addEventListener('orientationchange', adjustViewport);
adjustViewport(); // Initial call
