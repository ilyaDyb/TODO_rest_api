package ws

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/ilyaDyb/go_rest_api/logger"
	"github.com/sirupsen/logrus"
)

type Message struct {
	Username string `json:"username"`
	Content string `json:"content"`
}

type Client struct {
	ID string
	Conn *websocket.Conn
	Send chan Message
}

type Hub struct {
	Clients map[string]*Client
	Register chan *Client
	Unregister chan *Client
	Broadcast chan Message
	Mu sync.Mutex
}

var HubInstance = Hub{
	Clients: make(map[string]*Client),
	Register: make(chan *Client),
	Unregister: make(chan *Client),
	Broadcast: make(chan Message),
}

func NewClient(id string, conn *websocket.Conn) *Client {
	return &Client{
		ID: id,
		Conn: conn,
		Send: make(chan Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <- h.Register:
			h.Mu.Lock()
			h.Clients[client.ID] = client
			h.Mu.Unlock()
			log.Printf("Client connected: %s", client.ID)
			logger.Log.WithFields(logrus.Fields{
				"component": "websocket_chat",
			}).Infof("Client connected: %s", client.ID)

			case client := <- h.Unregister:
				h.Mu.Lock()
				if _, ok := h.Clients[client.ID]; ok {
					delete(h.Clients, client.ID)
					close(client.Send)
					log.Printf("Client disconnected: %s", client.ID)
					logger.Log.WithFields(logrus.Fields{
						"component": "websocket_chat",
					}).Infof("Client disconnected: %s", client.ID)
				}
			
			case message := <- h.Broadcast:
				h.Mu.Lock()
				for _, client := range h.Clients {
					select {
					case client.Send <- message:
					default:
						close(client.Send)
						delete(h.Clients, client.ID)
					}
				} 
				h.Mu.Unlock()
		}
	}
}


func (c *Client) ReadPump() {
    defer func() {
        log.Println("error when read pump")
        HubInstance.Unregister <- c
        c.Conn.Close()
    }()

    for {
        var msg Message
        if err := c.Conn.ReadJSON(&msg); err != nil {
            log.Printf("error reading message: %v", err.Error())
            logger.Log.WithFields(logrus.Fields{
                "component": "websocket_chat",
            }).Errorf("error reading message: %v", err.Error())
            break
        }
        log.Printf("client send message: %v", msg)
        HubInstance.Broadcast <- msg
        log.Println("message sent to Broadcast channel")
    }
}


func (c *Client) WritePump() {
    defer c.Conn.Close()
    for msg := range c.Send {
        log.Printf("server sending message: %v", msg)
        if err := c.Conn.WriteJSON(msg); err != nil {
            log.Printf("error writing message: %v", err.Error())
            logger.Log.WithFields(logrus.Fields{
                "component": "websocket_chat",
            }).Errorf("error writing message: %v", err.Error())
            break
        }
        log.Println("message sent to client")
    }
}