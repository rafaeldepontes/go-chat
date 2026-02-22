package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rafaeldepontes/go-chat/internal/middleware"
	"github.com/rafaeldepontes/go-chat/internal/tool"
)

func init() {
	envFile := ".env"
	tool.ChecksEnvFile(&envFile)
	godotenv.Load(envFile)
}

func main() {
	handler := middleware.NewHandler()

	go handler.MessageBroker.Process(handler.UserService.GetUserChannel())
	go handler.UserService.Save()

	fmt.Println("Server started running", os.Getenv("SERVER_URL"))
	log.Fatalln(handler.UserService.Run())
}
