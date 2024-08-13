package ws

import (
	"log"
	"net/http"

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
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.Error(c.Writer, "Failed to upgrade to websocket", http.StatusBadRequest)
		return
	}
	id := c.Query("id")
	client := NewClient(id, conn)
	log.Printf("Client was created with id: %v", id)
	HubInstance.Register <- client


	go client.ReadPump()
	go client.WritePump()
}