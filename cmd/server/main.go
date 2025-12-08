package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/rafaeldepontes/go-chat/internal/middleware"
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
	err := godotenv.Load(".env", ".env.example")
	if err != nil {
		fmt.Println("[ERROR] ", err)
		return
	}
	port := os.Getenv("SERVER_PORT")
	serverURL := os.Getenv("SERVER_URL")

	handler := middleware.NewHandler(&upgrader)

	go handler.Server.Run()

	fmt.Printf("API running on %v port\n", serverURL)
	http.ListenAndServe(":"+port, handler)
}
