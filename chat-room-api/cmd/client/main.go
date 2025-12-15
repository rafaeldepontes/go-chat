package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/rafaeldepontes/go-chat/internal/model"
	"github.com/rafaeldepontes/go-chat/internal/tool"
)

func clearLine() {
	fmt.Print("\x1b[1A\r\x1b[2K")
}

func read(conn *websocket.Conn) {
	var users []model.User
	var user model.User
	for {
		_, message, _ := conn.ReadMessage()

		if len(message) > 0 {
			if err := json.Unmarshal(message, &users); err == nil {
				for _, user := range users {
					fmt.Printf("%v: %v\n", user.Username, user.Message)
				}
			} else if err := json.Unmarshal(message, &user); err == nil {
				fmt.Printf("%v: %v\n", user.Username, user.Message)
			} else {
				fmt.Println("ERROR trying to deserialize the JSON:", err)
				continue
			}
		}
	}
}

var count = 1
var dialer = websocket.DefaultDialer

func main() {
	envFile := ".env"
	tool.ChecksEnvFile(&envFile)

	godotenv.Load(envFile)

	var serverURL string
	serverURL = os.Getenv("SERVER_URL")
	if os.Getenv("IS_TLS") != "false" {
		serverURL = os.Getenv("TLS_SERVER_URL")
	}

	certificate := os.Getenv("SERVER_CERTIFICATE")

	var user model.User
	user.Username = fmt.Sprintf("Anonymous%d", count)
	count++

	tlsConfig := &tls.Config{}

	if certificate != "" && tool.FileExists(certificate) && os.Getenv("IS_TLS") != "false" {
		caPEM, err := os.ReadFile(certificate)
		if err != nil {
			fmt.Println("failed to read CA cert:", err)
			return
		}

		rootPool := x509.NewCertPool()
		if !rootPool.AppendCertsFromPEM(caPEM) {
			fmt.Println("failed to append CA cert to pool")
			return
		}
		tlsConfig.RootCAs = rootPool
	}

	dialer.TLSClientConfig = tlsConfig
	conn, resp, err := dialer.Dial(serverURL, nil)
	if err != nil {
		if resp != nil {
			fmt.Printf("Dial error: %v (status: %s)\n", err, resp.Status)
			return
		}
		fmt.Printf("Dial error: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Print("Username: ")
	fmt.Scanln(&user.Username)
	fmt.Println("Connected!")
	fmt.Println("-------------------ChatRoom-------------------")

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
