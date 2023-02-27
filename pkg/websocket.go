package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// 验证来源是否允许
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan []byte)

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to upgrade WebSocket:", err)
		return
	}
	// 将每个客户端的 ID 存储到 clients map 中，以便在发送消息时可以识别客户端
	clients[conn] = true
	go handleMessages()
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseGoingAway) && !websocket.IsCloseError(err, websocket.CloseNoStatusReceived) {
				log.Println("Failed to read message from websocket:", err)
			}
			delete(clients, conn)
			return
		}
		log.Printf("Received message: %s\n", msg)
		if string(msg) == "sync_songs" {
			SyncSongs()
		}
	}

}

func handleMessages() {
	for {
		msg := <-broadcast
		for clientConn := range clients {
			err := clientConn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("Failed to send message to WebSocket:", err)
				delete(clients, clientConn)
			}
		}
	}
}

func SyncSongs() {
	opcode := "sync_songs "
	data, err := json.Marshal(GSync)
	if err != nil {
		log.Fatal(err)
	}
	broadcast <- []byte(opcode + string(data))
}

func SendMessage(ok int, msg string) {
	opcode := "send_msg "
	msgJson := ResultMessage{OK: ok, Msg: msg}
	data, err := json.Marshal(msgJson)
	if err != nil {
		log.Fatal(err)
	}
	broadcast <- []byte(opcode + string(data))
}
