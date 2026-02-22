package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rafaeldepontes/go-chat/internal/middleware"
	"github.com/rafaeldepontes/go-chat/internal/tool"

	messageapi "github.com/rafaeldepontes/go-chat/pkg/gRPC/message-api"
	"github.com/rafaeldepontes/go-chat/pkg/message-broker/rabbitmq"
)

func init() {
	envFile := ".env"
	tool.ChecksEnvFile(&envFile)

	err := godotenv.Load(envFile)
	if err != nil {
		fmt.Println("[ERROR] ", err)
		return
	}
}

func main() {
	port := os.Getenv("SERVER_PORT")

	var serverURL string
	getServerUrl(&serverURL)

	handler := middleware.NewHandler()

	go handler.Server.Run()
	defer rabbitmq.CloseConn()
	defer rabbitmq.CloseChan()
	defer messageapi.CloseConn()

	fmt.Printf("API running on %v port\n", serverURL)
	run(port, handler)
}

func run(port string, handler http.Handler) {
	if os.Getenv("IS_TLS") == "false" {
		log.Fatalln(http.ListenAndServe(":"+port, handler))
	} else {
		log.Fatalln(http.ListenAndServeTLS(
				":"+port,
				os.Getenv("SERVER_CERTIFICATE"),
				os.Getenv("SERVER_KEY"),
				handler,
			))
	}
}

func getServerUrl(serverURL *string) {
	*serverURL = os.Getenv("SERVER_URL")
	if os.Getenv("IS_TLS") != "false" {
		*serverURL = os.Getenv("TLS_SERVER_URL")
	}
}
