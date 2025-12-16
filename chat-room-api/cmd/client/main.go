package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"os"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/rafaeldepontes/go-chat/internal/model"
	"github.com/rafaeldepontes/go-chat/internal/tool"
)

var count = 1
var dialer = websocket.DefaultDialer

func main() {
	envFile := ".env"
	tool.ChecksEnvFile(&envFile)
	godotenv.Load(envFile)

	var serverURL string
	getServerUrl(&serverURL)

	var user model.User
	user.Username = fmt.Sprintf("Anonymous%d", count)
	count++

	tlsConfig := &tls.Config{}
	certificate := os.Getenv("SERVER_CERTIFICATE")
	useTLSProtection(certificate, tlsConfig)

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

	gettingUsername(&user)

	go read(conn)
	reader := bufio.NewReader(os.Stdin)

	run(reader, &user, conn)
}
