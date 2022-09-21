package service

import (
	"chatDemo/conf"
	"chatDemo/pkg/e"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

func (manager *ClientManager) Start() {
	for {
		log.Println("<---监听管道通信---->")
		select {
		case conn := <-manager.Register:
			log.Printf("建立新连接: %v", conn.ID)
			manager.Clients[conn.ID] = conn
			replyMsg := &ReplyMsg{
				Code:    e.WebsocketSuccess,
				Content: "已连接至服务器",
			}
			msg, err := json.Marshal(replyMsg)
			if err != nil {
				panic(err)
			}
			_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)

		case conn := <-manager.Unregister:
			log.Printf("连接失败: %v", conn.ID)
			if _, ok := manager.Clients[conn.ID]; ok {
				replyMsg := ReplyMsg{
					Code:    e.WebsocketEnd,
					Content: "连接已断开",
				}
				msg, err := json.Marshal(replyMsg)
				if err != nil {
					panic(err)
				}
				_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)
				close(conn.Send)
				delete(manager.Clients, conn.ID)
			}
		case broadCast := <-manager.Broadcast:
			message := broadCast.Message
			sendId := broadCast.Client.SendID
			flag := false // 默认对方不在线
			for id, conn := range manager.Clients {
				if id != sendId {
					continue
				}
				select {
				case conn.Send <- message:
					flag = true
				default:
					close(conn.Send)
					delete(manager.Clients, conn.ID)
				}

			}
			id := broadCast.Client.ID
			if flag {
				log.Println("对方在线应答")
				replyMsg := ReplyMsg{
					Code:    e.WebsocketOnlineReply,
					Content: "对方在线应答",
				}

				msg, err := json.Marshal(replyMsg)
				if err != nil {
					panic(err)
				}
				_ = broadCast.Client.Socket.WriteMessage(websocket.TextMessage, msg)
				err = InsertMsg(conf.MongoDBName, id, string(message), 1, int64(3*month))
				if err != nil {
					fmt.Println("InsertOneMsg Err", err)
				}

			} else {
				log.Println("对方不在线")
				replyMsg := ReplyMsg{
					Code:    e.WebsocketOfflineReply,
					Content: "对方不在线应答",
				}
				msg, err := json.Marshal(replyMsg)
				if err != nil {
					panic(err)
				}
				_ = broadCast.Client.Socket.WriteMessage(websocket.TextMessage, msg)
				err = InsertMsg(conf.MongoDBName, id, string(message), 0, int64(3*month))
				if err != nil {
					fmt.Println("InsertOneMsg Err", err)
				}
			}

		}

	}
}
