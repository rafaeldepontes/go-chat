# Go-Chat

[![language](https://img.shields.io/badge/language-Go-00ADD8?labelColor=2F2F2F)](https://go.dev/doc/)
[![version](https://img.shields.io/badge/version-1.25-9C27B0?labelColor=2F2F2F)](https://go.dev/doc/install)

<!-- [![build](https://img.shields.io/github/actions/workflow/status/rafaeldepontes/go-chat/build.yml?label=build&logo=githubactions&logoColor=white&labelColor=2F2F2F)](https://github.com/rafaeldepontes/go-chat/actions/workflows/build.yml)
[![tests](https://img.shields.io/github/actions/workflow/status/rafaeldepontes/go-chat/tests.yml?label=tests&logo=go&logoColor=white&labelColor=2F2F2F)](https://github.com/rafaeldepontes/go-chat/actions/workflows/tests.yml) -->

---

This repository is a hands-on technical demonstration and learning tool built in Go that focuses on real-time communication using WebSockets.

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

Go-Chat is a small WebSocket API that demonstrates server/client communication patterns in Go. It is intended as a learning project to test connection handling, message routing, and simple client identification.

## Features

- Minimal WebSocket server implementation
- Simple client that connects to the server and exchanges messages
- Example of environment-based configuration
- Small, easy-to-read codebase intended for experimentation and extension

## Technologies

- Go 1.25
- `gorilla/websocket`
- `godotenv`

## Requirements

- Go 1.25 or later

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

## Environment variables

Copy `.env.example` to `.env` and update values as needed. The project expects at least the following variables:

```bash
SERVER_PORT="8000"
SERVER_URL="ws://localhost:8000"

# ----------------------------

POSTGRES_DB="chatroom-db"
POSTGRES_USER="root"
POSTGRES_PASSWORD="example"
DATABASE_URL="postgres://root:example@localhost:5432/postgres"
```

## Quick start

1. Clone the repository and enter the folder:

   ```bash
   git clone <repo-url>
   cd go-chat
   ```

2. Copy the example environment file and adjust values:

   ```bash
   cp .env.example .env
   # or create .env manually and set SERVER_PORT and SERVER_URL
   ```

3. Run the server first, then one or more clients:

   ```bash
   go run cmd/server/main.go   # start server
   go run cmd/client/main.go   # start a client (you can run many clients)
   ```

4. In the client you will be prompted for a username. After connecting you can type messages that will be sent to the server and routed accordingly.

## Usage and behavior notes

- The server binds to the port defined in `SERVER_PORT` and expects clients to connect to `SERVER_URL`.
- Clients are lightweight examples meant to demonstrate how to connect, send, and receive messages.
- Clients run on random local ports and are intended for local testing only.
- On entry a client should see all the messages sended in that chat.
- `Cache` stores recent messages to avoid repeated DB hits, it uses TTL.

### Message Flow

1. User sends a message
2. Server persists it to PostgreSQL
3. Server updates the in‑memory cache
4. Server broadcasts message to all connected clients

## Development notes

- The codebase is intentionally simple so you can experiment with connection handling, message formats, and routing strategies.
- Ideas to extend the project:

  - Add tests covering core connection and routing logic `WIP`

## License

[Click here](LICENSE)

## Contact

For questions or help, contact: `rafael.cr.carneiro@gmail.com`
