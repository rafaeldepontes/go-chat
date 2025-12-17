# Go-Chat

[![language](https://img.shields.io/badge/language-Go-00ADD8?labelColor=2F2F2F)](https://go.dev/doc/)
[![version](https://img.shields.io/badge/version-1.25-9C27B0?labelColor=2F2F2F)](https://go.dev/doc/install)

<!-- [![build](https://img.shields.io/github/actions/workflow/status/rafaeldepontes/go-chat/build.yml?label=build&logo=githubactions&logoColor=white&labelColor=2F2F2F)](https://github.com/rafaeldepontes/go-chat/actions/workflows/build.yml)
[![tests](https://img.shields.io/github/actions/workflow/status/rafaeldepontes/go-chat/tests.yml?label=tests&logo=go&logoColor=white&labelColor=2F2F2F)](https://github.com/rafaeldepontes/go-chat/actions/workflows/tests.yml) -->

---

This repository is a hands-on technical demonstration and learning tool built in Go that focuses on real-time communication using WebSockets, a asyncronous saving using RabbitMQ and Postgres for eventual consistency.

## Table of contents

- [Overview](#overview)
- [Features](#features)
- [Technologies](#technologies)
- [Requirements](#requirements)
- [Environment variables](#environment-variables)
- [Quick start](#quick-start)
- [Usage](#usage-and-behavior-notes)
- [Message Flow](#message-flow)
- [Development notes](#development-notes)
- [License](#license)
- [Contact](#contact)

## Overview

Go-Chat is a small WebSocket API that demonstrates server/client communication patterns in Go. It is intended as a learning project to test connection handling, message routing, and simple client identification, RabbitMQ asynchronous calls and gRPC connections and uses.

## Features

- Minimal WebSocket server implementation
- Simple client that connects to the server and exchanges messages
- Example of environment-based configuration
- Small, easy-to-read codebase intended for experimentation and extension
- RabbitMQ for eventual consistent
- gRPC for server-to-server comunication

## Technologies

- Go 1.25
- `gorilla/websocket`
- `rabbitmq/amqp0910-go`
- `google.golang.org/grpc`
- `godotenv`

## Requirements

- Go 1.25 or later

## Environment variables

Copy `.env.example` to `.env` from both folders, `chat-room-api` and `message-api` and update values as needed. The project expects at least the following variables:

## Chat-Room-API:
```bash
# Change this as needed, keep in mind that every request without
# the certificate and with this variable changed can crash the
# entire system
IS_TLS="false"
# IS_TLS="true"

# ----------------------------

SERVER_PORT="8080"
SERVER_URL="ws://localhost:8080"
TLS_SERVER_URL="wss://localhost:8080"

# ----------------------------

SERVER_KEY="server.key"
SERVER_CERTIFICATE="server.crt"

# ----------------------------

RABBITMQ_URL="amqp://user:pass@localhost:5672/"

# ----------------------------

MESSAGE_SERVICE_PORT="localhost:8081"
```

## Message-API:
```bash
# Change this as needed, keep in mind that every request without
# the certificate and with this variable changed can crash the
# entire system
IS_TLS="false"
# IS_TLS="true"

# ----------------------------

SERVER_PORT="8081"
SERVER_URL="http://localhost:8081"
TLS_SERVER_URL="https://localhost:8081"

# ----------------------------

POSTGRES_DB="chatroom-db"
POSTGRES_USER="root"
POSTGRES_PASSWORD="example"
DATABASE_URL="postgres://root:example@localhost:5432/postgres"

# ----------------------------

SERVER_KEY="server.key"
SERVER_CERTIFICATE="server.crt"

# ----------------------------

RABBITMQ_URL="amqp://user:pass@localhost:5672/"

```

## Database Schema

Use the schema below to initialize the database:

```sql
create table chat_room (
	id serial primary key,
	username varchar(50),
	message varchar(512),
	sent_at TIME default now()
);
```

## Quick start

1. Clone the repository and enter the folder:

   ```bash
   git clone <repo-url>

   cd ./go-chat/message-api/
   go mod tidy

   cd ../chat-room-api/
   go mod tidy
   ```

2. Copy the example environment file and adjust values:

   ```bash
   cp .env.example .env
   # or create .env manually and set SERVER_PORT and SERVER_URL
   ```

3. In the root folder, starts both PostgreSQL and RabbitMQ with Docker:

   ```bash
   docker-compose up -d
   ```

4. Apply the database schema (run the SQL above) using `psql` or a GUI tool.

### OPTIONAL

5. If you want to use a HTTPS connection, you need to change the `.env` or the `.env.example` file in such manner:

   ```bash
   # IS_TLS="false"  # <-- comment this section.
   IS_TLS="true"
   ```

- 5.1 When you enable `TLS`, generating a certificate is mandatory. Run this command in your `"bash" (exclusive to Linux)`:

  ```bash
  openssl req -x509 -nodes -newkey rsa:2048 -keyout server.key -out server.crt -days 3650 -subj "//CN=localhost" -addext "subjectAltName = DNS:localhost,IP:127.0.0.1,IP:::1"
  ```

---

6. In the client you will be prompted for a username. After connecting you can type messages that will be sent to the server and routed accordingly.

7. To run this application you need to have at least three cli's opened at the same time and both the database and rabbitmq running...

   7.1 Inside the root folder (./), run the following commands:
   ```bash
   cd ./chat-room-api/
   go run cmd/server/main.go
   ```

   7.2 Inside the root folder (./), run the following commands:
   ```bash
   cd ./message-api/
   go run .
   ```

   7.3 Inside the root folder (./), run the following commands:
   ```bash
   cd ./chat-room-api/
   go run cmd/client/*.go # You can have as many as you want running...
   ```

## Usage and behavior notes

- The server binds to the port defined in `SERVER_PORT` and expects clients to connect to `SERVER_URL`.
- Clients are lightweight examples meant to demonstrate how to connect, send, and receive messages.
- Clients run on random local ports and are intended for local testing only.
- On entry a client should see all the messages sended in that chat.
- `Cache` stores recent messages to avoid repeated DB hits, it uses TTL. `Deprecated on migration, I'm planning to create it again`

### Message Flow

1. User chooses a name
2. If there's any message stored in the database, the chat room server will make a gRPC call to the message-api and it should list every single message in the table.
3. User/Client sends a message
4. Chat Room Server sends the message to RabbitMQ
5-1. Server broadcasts message to all connected clients
5-2. Message API consumes the queue and persists it to PostgreSQL

## Development notes

- The codebase is intentionally simple so you can experiment with connection handling, message formats, and routing strategies.
- Ideas to extend the project:

  - Add tests covering core connection and routing logic `WIP`

## License

[Click here](LICENSE)

## Contact

For questions or help, contact: `rafael.cr.carneiro@gmail.com`
