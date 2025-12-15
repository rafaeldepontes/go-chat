package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rafaeldepontes/go-chat/internal/middleware"
	"github.com/rafaeldepontes/go-chat/internal/tool"
)

func main() {
	envFile := ".env"
	tool.ChecksEnvFile(&envFile)
	godotenv.Load(envFile)

	var port string
	port = ":" + os.Getenv("SERVER_PORT")

	handler := middleware.NewHandler()

	go handler.MessageBroker.Process(handler.UserService.GetUserChannel())
	go handler.UserService.Save()

	fmt.Println("Server started running", os.Getenv("SERVER_URL"))

	if err := http.ListenAndServe(port, handler); err != nil {
		fmt.Println("[ERROR]", err)
	}
}
