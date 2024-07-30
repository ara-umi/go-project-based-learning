package main

import (
	"encoding/json"
	"fmt"
)

type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func (manager *ClientManager) start() {
	for {
		select {
		// 注册 chan 来消息（消息都是 pclient），就写进 clients 这个 map
		case conn := <-manager.register:
			manager.clients[conn] = true
			jsonMessage, _ := json.Marshal(
				&Message{
					Sender:  "admin",
					Content: fmt.Sprintf("%s connected", conn.id),
				},
			)
			manager.send(jsonMessage, conn) // 实际上功能类似广播，通知所有用户，有人连进来了（除了他自己）
		case conn := <-manager.unregister:
			if _, ok := manager.clients[conn]; ok { // 确认已注册
				close(conn.send)              // 移除逻辑1
				delete(manager.clients, conn) // 移除逻辑2
				jsonMessage, _ := json.Marshal(
					&Message{
						Sender:  "admin",
						Content: fmt.Sprintf("%s disconnected", conn.id),
					},
				)
				manager.send(jsonMessage, conn) // 通知所有用户有人断连了（除了他自己，当然也没必要，因为已经 delete 了）
			}
		case message := <-manager.broadcast:
			// 广播 chan 来消息，就对 map 中的所有 client 的 send chan []byte 尝试写入 message
			for conn := range manager.clients {
				select {
				case conn.send <- message: // 非阻塞，能写入，后续也不需要执行什么逻辑
				default: // 阻塞，就当断开连接，执行移除的逻辑 12
					close(conn.send)
					delete(manager.clients, conn)
				}
			}
		}
	}
}

// 对 manager 中管理的所有不是 ignore 的 client 的 send 这个 chan []byte 发送消息
func (manager *ClientManager) send(message []byte, ignore *Client) {
	for conn := range manager.clients {
		if conn != ignore {
			select {
			case conn.send <- message:
			default:
				close(conn.send)
				delete(manager.clients, conn)
			}
		}
	}
}
