package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

type User struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

func clearLine() {
	fmt.Print("\x1b[1A\r\x1b[2K")
}

func read(conn *websocket.Conn) {
	var user User
	for {
		_, message, _ := conn.ReadMessage()
		json.Unmarshal(message, &user)
		fmt.Printf("%v: %v\n", user.Username, user.Message)
	}
}

func main() {
	godotenv.Load(".env", ".env.example")
	serverURL := os.Getenv("SERVER_URL")

	var user User

	fmt.Print("Username: ")
	fmt.Scanln(&user.Username)
	fmt.Println("-------------------ChatRoom-------------------")

	conn, resp, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		if resp != nil {
			fmt.Printf("Dial error: %v (status: %s)\n", err, resp.Status)
			return
		}
		fmt.Printf("Dial error: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected!")

	go read(conn)
	reader := bufio.NewReader(os.Stdin)

	for {
		input, _ := reader.ReadString('\n')
		clearLine()
		text := strings.TrimSpace(input)
		if text == "" {
			continue
		}

		user.Message = text

		fmt.Printf("%v: %v\n", user.Username, user.Message)

		jsonBody, _ := json.Marshal(user)
		err = conn.WriteMessage(websocket.TextMessage, jsonBody)
		if err != nil {
			fmt.Printf("Error writing: %v\n", err)
			return
		}
	}
}
