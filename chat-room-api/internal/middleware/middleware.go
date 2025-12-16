package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/rafaeldepontes/go-chat/internal/server/service"
)

type handler struct {
	Upgrader *websocket.Upgrader
	Server   *service.Server
}

const (
	ReadBufferSize  = 4096
	WriteBufferSize = 4096
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  ReadBufferSize,
	WriteBufferSize: WriteBufferSize,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func NewHandler() *handler {
	return &handler{
		Upgrader: &upgrader,
		Server:   service.NewService(),
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var conn *websocket.Conn
	conn, err := h.Upgrader.Upgrade(w, r, nil)
	if err != nil && !errors.Is(err, http.ErrHijacked) {
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
