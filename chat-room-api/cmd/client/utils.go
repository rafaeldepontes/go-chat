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
	"github.com/rafaeldepontes/go-chat/internal/model"
	"github.com/rafaeldepontes/go-chat/internal/tool"
)

func getServerUrl(serverURL *string) {
	*serverURL = os.Getenv("SERVER_URL")
	if os.Getenv("IS_TLS") != "false" {
		*serverURL = os.Getenv("TLS_SERVER_URL")
	}
}

func useTLSProtection(certificate string, tlsConfig *tls.Config) {
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
}

func gettingUsername(user *model.User) {
	fmt.Print("Username: ")
	fmt.Scanln(&user.Username)
	fmt.Println("Connected!")
	fmt.Println("-------------------ChatRoom-------------------")
}

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

func run(reader *bufio.Reader, user *model.User, conn *websocket.Conn) {
	for {
		input, _ := reader.ReadString('\n')
		clearLine()
		text := strings.TrimSpace(input)
		if text == "" {
			continue
		}

		user.Message = text

		fmt.Printf("%v: %v\n", user.Username, user.Message)

		jsonBody, _ := json.Marshal(*user)
		err := conn.WriteMessage(websocket.TextMessage, jsonBody)
		if err != nil {
			fmt.Printf("Error writing: %v\n", err)
		}
	}
}
