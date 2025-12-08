package service

import "sync"

type Server struct {
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
	mux        sync.Mutex
}

func NewService() *Server {
	return &Server{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
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

		case client := <-s.Unregister:
			delete(s.Clients, client)
			close(client.Send)
			_ = client.Conn.Close()
		}
	}
}
