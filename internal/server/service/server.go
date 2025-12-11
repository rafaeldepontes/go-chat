package service

import (
	"fmt"
	"sync"

	"github.com/rafaeldepontes/go-chat/internal/user"
	"github.com/rafaeldepontes/go-chat/internal/user/service"
)

type Server struct {
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
	UserSvc    user.Service
	mux        sync.Mutex
}

func NewService() *Server {
	return &Server{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
		UserSvc:    service.NewService(),
		mux:        sync.Mutex{},
	}
}

func (s *Server) Run() {
	var c *Client
	for {
		select {
		case msg := <-s.Broadcast:
			for client := range s.Clients {
				if !client.isSending {
					client.Send <- msg
				} else {
					c = client
				}
			}

			c.isSending = false

		case client := <-s.Register:
			s.Clients[client] = true

			msg, err := s.UserSvc.FindAll()
			if err != nil {
				fmt.Println("An unexpected error while reading the database:", err)
				continue
			}

			client.Send <- msg

		case client := <-s.Unregister:
			delete(s.Clients, client)
			close(client.Send)
			_ = client.Conn.Close()
		}
	}
}
