package service

import (
	"bytes"
	"fmt"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn      *websocket.Conn
	Send      chan []byte
	isSending bool
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		Conn:      conn,
		Send:      make(chan []byte),
		isSending: false,
	}
}

func (c *Client) Read(s *Server) {
	defer func() {
		s.Unregister <- c
		_ = c.Conn.Close()
	}()

	for {
		mt, message, err := c.Conn.ReadMessage()
		if err != nil {
			fmt.Println("Error while reading the message:", err)
			return
		}

		if mt == websocket.BinaryMessage {
			fmt.Println("Not supported binary")
			return
		}

		s.mux.Lock()
		c.isSending = true
		s.mux.Unlock()

		cleanMessage := bytes.Trim(message, "/n")
		s.Broadcast <- cleanMessage
	}
}

func (c *Client) SendMessage() {
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte("Connection closed"))
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				fmt.Println("An error occurred:", err)
				fmt.Println("Cannot write message, connection closed...")
				return
			}
		default:
		}
	}
}
