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

		cleanMessage := bytes.Trim(message, "/n")

		s.mux.Lock()
		c.isSending = true
		if err = s.UserSvc.Save(cleanMessage...); err != nil {
			fmt.Println("Error while trying to save the message:", err)
			return
		}
		s.mux.Unlock()

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
