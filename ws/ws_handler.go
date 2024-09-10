package ws

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)


var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}


func WsHandler(c *gin.Context) {
	chatID := c.Param("chatID")
	username := c.Param("username")

	chatIDInt, err := strconv.Atoi(chatID)
	if err != nil {
		http.Error(c.Writer, "Invalid chat ID", http.StatusBadRequest)
		return
	}
	chatIDUint := uint(chatIDInt)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.Error(c.Writer, "Failed to upgrade to websocket", http.StatusBadRequest)
		return
	}

	client := NewClient(chatIDUint, username, conn)
	log.Printf("Client was created with username: %v and chatID: %v", username, chatIDUint)

	if len(HubInstance.Chats[chatIDUint]) >= 2 {
		client.Conn.WriteMessage(websocket.TextMessage, []byte("Chat is full"))
		client.Conn.Close()
		return
	}

	HubInstance.Register <- client

	go client.ReadPump()
	go client.WritePump()
}