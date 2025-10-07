package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

func (c *Client) ReadPump() {
	defer func() {
		HubInstance.Unregister <- c
		c.Conn.Close()
	}()
	for {
		if _, _, err := c.Conn.ReadMessage(); err != nil {
			log.Println("read error:", err)
			break
		}
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()
	for msg := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("write error:", err)
			break
		}
	}
}
