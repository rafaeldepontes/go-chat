package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/rafaeldepontes/go-chat/internal/middleware"
	"github.com/rafaeldepontes/go-chat/internal/tool"
	// "github.com/rafaeldepontes/go-chat/pkg/db/postgres"
	"github.com/rafaeldepontes/go-chat/pkg/message-broker/rabbitmq"
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

	var serverURL string
	serverURL = os.Getenv("SERVER_URL")
	if os.Getenv("IS_TLS") != "false" {
		serverURL = os.Getenv("TLS_SERVER_URL")
	}

	handler := middleware.NewHandler(&upgrader)

	go handler.Server.Run()
	// defer postgres.Disconnect()
	defer rabbitmq.CloseConn()
	defer rabbitmq.CloseChan()

	fmt.Printf("API running on %v port\n", serverURL)
	if os.Getenv("IS_TLS") != "false" {
		if err := http.ListenAndServeTLS(":"+port, os.Getenv("SERVER_CERTIFICATE"), os.Getenv("SERVER_KEY"), handler); err != nil {
			fmt.Println("Error initializing the api:", err)
		}
	} else {
		if err := http.ListenAndServe(":"+port, handler); err != nil {
			fmt.Println("Error initializing the api:", err)
		}
	}
}
