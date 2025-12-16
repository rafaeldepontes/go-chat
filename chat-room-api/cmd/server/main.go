package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rafaeldepontes/go-chat/internal/middleware"
	"github.com/rafaeldepontes/go-chat/internal/tool"

	"github.com/rafaeldepontes/go-chat/pkg/message-broker/rabbitmq"
)

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
	getServerUrl(&serverURL)

	handler := middleware.NewHandler()

	go handler.Server.Run()
	defer rabbitmq.CloseConn()
	defer rabbitmq.CloseChan()

	fmt.Printf("API running on %v port\n", serverURL)
	run(port, handler)
}

func run(port string, handler http.Handler) {
	if os.Getenv("IS_TLS") == "false" {
		if err := http.ListenAndServe(":"+port, handler); err != nil {
			fmt.Println("Error initializing the api:", err)
		}
	} else {
		if err := http.ListenAndServeTLS(":"+port, os.Getenv("SERVER_CERTIFICATE"), os.Getenv("SERVER_KEY"), handler); err != nil {
			fmt.Println("Error initializing the api:", err)
		}
	}
}

func getServerUrl(serverURL *string) {
	*serverURL = os.Getenv("SERVER_URL")
	if os.Getenv("IS_TLS") != "false" {
		*serverURL = os.Getenv("TLS_SERVER_URL")
	}
}
