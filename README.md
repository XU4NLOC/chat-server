# chat-server

A real-time chat server built with Go, featuring WebSocket-based messaging, JWT authentication, and chat rooms. Includes a clean web-based UI.

## Features

- **Real-time messaging** via WebSocket connections
- **User authentication** with JWT tokens (register & login)
- **Chat rooms** — create, list, and join rooms
- **Message history** — last 50 messages loaded on join
- **Online presence** — live user count per room
- **Live room creation** — new rooms broadcast to all connected clients
- **Session persistence** — token stored in `localStorage` for auto-restore
- **Docker support** — multi-stage build + `docker-compose`

## Tech Stack

| Layer       | Technology                          |
|-------------|-------------------------------------|
| Language    | Go 1.26                             |
| WebSocket   | [gorilla/websocket](https://github.com/gorilla/websocket) |
| Auth        | [golang-jwt](https://github.com/golang-jwt/jwt) (HS256) |
| Password    | bcrypt (via `golang.org/x/crypto`)  |
| Database    | PostgreSQL 16 (via `lib/pq`)        |
| Frontend    | Vanilla HTML/CSS/JS                 |
| Infra       | Docker, docker-compose              |

## Project Structure

```
chat-server/
├── main.go              # Entry point — wires routes, hub, and DB
├── config/config.go     # Environment variable loading
├── db/db.go             # PostgreSQL connection & auto-migration
├── auth/jwt.go          # JWT token generation & validation
├── models/
│   ├── user.go          # User CRUD, bcrypt password handling
│   ├── room.go          # Room CRUD
│   └── message.go       # Message queries
├── handlers/
│   ├── auth.go          # POST /register, POST /login
│   ├── room.go          # GET /rooms, POST /rooms
│   └── ws.go            # WebSocket upgrade & client setup
├── hub/
│   ├── hub.go           # Central broadcast hub (room-based routing)
│   └── client.go        # WebSocket read/write pumps per client
├── static/index.html    # Web-based chat UI
├── Dockerfile           # Multi-stage build
├── docker-compose.yml   # App + PostgreSQL orchestration
└── go.mod
```

## Getting Started

### Prerequisites

- Go 1.26+
- PostgreSQL 16+ (or use Docker)
- Docker & Docker Compose (optional)

### Option 1: Docker Compose (Recommended)

```bash
docker compose up --build
```

The app will be available at `http://localhost:8080`.

This starts both the Go server and a PostgreSQL instance. The database is auto-migrated on startup.

### Option 2: Run Locally

1. Start a PostgreSQL instance and create a database:

```bash
createdb chatdb
```

1. Set environment variables (or use the defaults):

```bash
export DATABASE_URL="postgres://localhost:5432/chatdb?sslmode=disable"
export JWT_SECRET="your-secret-key"
export PORT="8080"
```

1. Build and run:

```bash
go build -o chat-server .
./chat-server
```

The server starts on `http://localhost:8080`. Database tables are created automatically on first run.

## API Reference

### Authentication

| Method | Endpoint     | Body                              | Description          |
|--------|--------------|-----------------------------------|----------------------|
| POST   | `/register`  | `{"username": "...", "password": "..."}` | Create a new account |
| POST   | `/login`     | `{"username": "...", "password": "..."}` | Get a JWT token      |

**Login response:**

```json
{ "token": "eyJhbGciOi..." }
```

### Rooms

| Method | Endpoint  | Body                      | Description              |
|--------|-----------|---------------------------|--------------------------|
| GET    | `/rooms`  | —                         | List all rooms           |
| POST   | `/rooms`  | `{"name": "room-name"}`   | Create a new room        |

### WebSocket

Connect to:

```
ws://localhost:8080/ws?token=<JWT>&room_id=<ID>
```

**Query parameters:**

- `token` — JWT from `/login`
- `room_id` — Room to join

**Incoming message format:**

```json
{ "username": "alice", "content": "hello!", "room_id": 1 }
```

**System events:**

```json
{ "type": "system", "event": "user_count", "count": 3 }
{ "type": "system", "event": "new_room", "room": { "id": 4, "name": "general" } }
```

**Outgoing message format (send from client):**

```json
{ "content": "hello!" }
```

## Environment Variables

| Variable      | Default                                                       | Description            |
|---------------|---------------------------------------------------------------|------------------------|
| `DATABASE_URL`| `postgres://lexuanloc@localhost:5432/chatdb?sslmode=disable`  | PostgreSQL connection  |
| `JWT_SECRET`  | `supersecretkey`                                              | HMAC signing secret    |
| `PORT`        | `8080`                                                        | Server listen port     |

## Docker Deployment

Build for a specific platform (e.g., for cloud deployment):

```bash
docker build --platform=linux/amd64 -t chat-server .
docker push your-registry/chat-server
```

See [README.Docker.md](README.Docker.md) for more details.

## License

MIT

