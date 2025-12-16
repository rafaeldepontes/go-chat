package service

import (
	"fmt"

	"github.com/rafaeldepontes/go-chat/internal/message"
	messagebroker "github.com/rafaeldepontes/go-chat/internal/message-broker"
	msgBrokerSvc "github.com/rafaeldepontes/go-chat/internal/message-broker/service"
	msgSvc "github.com/rafaeldepontes/go-chat/internal/message/service"
)

type Server struct {
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
	MessageSvc message.Service
	MsgBroker  messagebroker.MsgBroker
}

func NewService() *Server {
	return &Server{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
		MessageSvc: msgSvc.NewService(),
		MsgBroker:  msgBrokerSvc.NewMsgService(),
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

			msg, err := s.MessageSvc.FindAll()
			if err != nil {
				fmt.Println("An unexpected error happened while trying to fetch the chat data...\nError:", err)
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
