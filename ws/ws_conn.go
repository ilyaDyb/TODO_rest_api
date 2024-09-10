package ws

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/ilyaDyb/go_rest_api/config"
	"github.com/ilyaDyb/go_rest_api/logger"
	"github.com/ilyaDyb/go_rest_api/models"
	"github.com/sirupsen/logrus"
)

type Message struct {
	ChatID     uint   `json:"chat_id"`
	SenderID   uint   `json:"sender_id"`
	ReceiverID uint   `json:"receiver_id"`
	Content  string `json:"content"`
}

type Client struct {
	ID    string
	Conn  *websocket.Conn
	Send  chan Message
	ChatID uint
}

type Hub struct {
	Chats     map[uint]map[string]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message
	Mu         sync.Mutex
}

var HubInstance = &Hub{
	Chats:      make(map[uint]map[string]*Client),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
	Broadcast:  make(chan Message),
}

func NewClient(chatID uint, username string, conn *websocket.Conn) *Client {
	return &Client{
		ID:     username,
		Conn:   conn,
		Send:   make(chan Message),
		ChatID: chatID,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Mu.Lock()
			chat := client.ChatID
			
			if h.Chats[chat] == nil {
				h.Chats[chat] = make(map[string]*Client)
			}

			h.Chats[chat][client.ID] = client
			// log.Printf("Client %s joined chat %s", client.ID, chat)

			h.Mu.Unlock()
			log.Printf("Client connected: %s", client.ID)
			logger.Log.WithFields(logrus.Fields{
				"component": "websocket_chat",
			}).Infof("Client connected: %s", client.ID)

		case client := <-h.Unregister:
			h.Mu.Lock()
			chat := client.ChatID

			if clients, ok := h.Chats[chat]; ok {
				if _, exists := clients[client.ID]; exists {
					delete(clients, client.ID)
					close(client.Send)
					log.Printf("Client disconnected: %s", client.ID)
					logger.Log.WithFields(logrus.Fields{
						"component": "websocket_chat",
					}).Infof("Client disconnected: %s", client.ID)

					if len(clients) == 0 {
						delete(h.Chats, chat)
					}
				}
			}
			h.Mu.Unlock()

		case message := <-h.Broadcast:
			h.Mu.Lock()
			chat := message.ChatID
			if clients, ok := h.Chats[chat]; ok {
				for _, client := range clients {
					select {
					case client.Send <- message:
					default:
						close(client.Send)
						delete(clients, client.ID)
					}
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
		msg.ChatID = c.ChatID
		log.Printf("client send message: %v", msg)

		isRead := len(HubInstance.Chats[c.ChatID]) == 2
		if err := createMessage(&msg, isRead); err != nil {
			HubInstance.Unregister <- c
			c.Conn.Close()
			break
		}
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

func createMessage(msg *Message, isRead bool) error {
	message := models.Message{
		ChatID: msg.ChatID,
		SenderID: msg.SenderID,
		ReceiverID: msg.ReceiverID,
		Content: msg.Content,
	}

	if err := config.DB.Create(&message).Error; err != nil {
		return err
	}
	// log.Println("Message was created successfully", message)
	return nil
}