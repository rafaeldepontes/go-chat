package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/rafaeldepontes/go-chat/internal/middleware"
	"github.com/rafaeldepontes/go-chat/internal/tool"
	"github.com/rafaeldepontes/go-chat/pkg/db/postgres"
)

const (
	ReadBufferSize  = 4096
	WriteBufferSize = 4096
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  ReadBufferSize,
	WriteBufferSize: WriteBufferSize,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func main() {
	envFile := ".env"
	tool.ChecksEnvFile(&envFile)

	err := godotenv.Load(envFile)
	if err != nil {
		fmt.Println("[ERROR] ", err)
		return
	}
	port := os.Getenv("SERVER_PORT")
	serverURL := os.Getenv("SERVER_URL")

	handler := middleware.NewHandler(&upgrader)

	go handler.Server.Run()
	defer postgres.Disconnect()

	fmt.Printf("API running on %v port\n", serverURL)
	http.ListenAndServe(":"+port, handler)
}
