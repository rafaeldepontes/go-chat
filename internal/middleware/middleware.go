package middleware

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/rafaeldepontes/go-chat/internal/server/service"
)

type handler struct {
	Upgrader *websocket.Upgrader
	Server   *service.Server
}

func NewHandler(upgrader *websocket.Upgrader) *handler {
	return &handler{
		Upgrader: upgrader,
		Server:   service.NewService(),
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var conn *websocket.Conn
	conn, err := h.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error while trying to establish a connection")
		fmt.Println("An error occurred:", err)
		return
	}
	fmt.Println("New connection established!")

	client := service.NewClient(conn)
	h.Server.Register <- client

	go client.Read(h.Server)
	go client.SendMessage()

	h.ServeHTTP(w, r)
}
