package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/rafaeldepontes/go-chat/internal/middleware"
	"github.com/rafaeldepontes/go-chat/internal/tool"
)

func main() {
	envFile := ".env"
	tool.ChecksEnvFile(&envFile)
	godotenv.Load(envFile)

	handler := middleware.NewHandler()

	go handler.MessageBroker.Process(handler.UserService.GetUserChannel())
	go handler.UserService.Save()

	fmt.Println("Server started running", os.Getenv("SERVER_URL"))

	if err := handler.UserService.Run(); err != nil {
		fmt.Println("[ERROR] while server initialization...\n[Error]:", err)
	}
}
