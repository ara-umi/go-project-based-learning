package main

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type Client struct {
	id     string
	socket *websocket.Conn
	send   chan []byte
}

func (c *Client) read() {
	defer func() {
		Manager.unregister <- c
		c.socket.Close()
	}()

	// 主逻辑
	for {
		_, message, err := c.socket.ReadMessage() // 从 ws 读消息
		if err != nil {                           // 逻辑同 defer，读错误，就断连并退出主逻辑
			Manager.unregister <- c
			c.socket.Close()
			break
		}
		jsonMessage, _ := json.Marshal(&Message{Sender: c.id, Content: string(message)})
		Manager.broadcast <- jsonMessage // 读的消息发送到 manager 的 broadcast
		// Logger.Infof("Message received: %s", string(message))
	}
}

func (c *Client) write() {
	defer func() {
		c.socket.Close()
	}()

	for message := range c.send {
		c.socket.WriteMessage(websocket.TextMessage, message)
		// Logger.Infof("Message sent: %s", string(message))
	}
}
