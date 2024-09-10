package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ilyaDyb/go_rest_api/config/redis"
	"github.com/ilyaDyb/go_rest_api/logger"
	"github.com/ilyaDyb/go_rest_api/models"
	"github.com/ilyaDyb/go_rest_api/service"
	"github.com/ilyaDyb/go_rest_api/tasks"
	"github.com/sirupsen/logrus"
)

type ChatController struct {
	chatService service.ChatService
	userService service.UserService
}

func NewChatController(chatService service.ChatService, userService service.UserService) *ChatController {
	return &ChatController{
		chatService: chatService,
		userService: userService,
	}
}

// func (ctrl *ChatController) CreateChat(c *gin.Context) {
// 	var input struct {
// 		User1ID uint `json:"user1_id"`
// 		User2ID uint `json:"user2_id"`
// 	}

// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		logger.Log.WithFields(logrus.Fields{
// 			"component": "chat",
// 		}).Errorf("client sent invalid data with error: %v", err.Error())
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	chat, err := ctrl.chatService.CreateChat(input.User1ID, input.User2ID)
// 	if err != nil {
// 		logger.Log.WithFields(logrus.Fields{
// 			"component": "chat",
// 		}).Errorf("could not create chat with error: %v", err.Error())
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create chat"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, chat)
// }

// ChatController godoc
// @Summary Get a chat between two users
// @Description Get a chat between the current user and the target user
// @Tags chat
// @Produce json
// @Param Authorization header string true "With the Bearer started"
// @Param username path string true "Target Username"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /chats/{username} [get]
func (ctrl *ChatController) ChatController(c *gin.Context) {
    targetUsername := c.Param("username")
    if targetUsername == "" {
        logger.Log.WithFields(logrus.Fields{
            "component": "chat",
        }).Info("route must contain target username")
        c.JSON(http.StatusBadRequest, gin.H{"error": "route must contain target username"})
        return
    }

    username := c.MustGet("username").(string)
    if username == targetUsername {
        c.JSON(http.StatusBadRequest, gin.H{"error": "you can't join a chat with yourself"})
        return
    }

    chat, err := ctrl.chatService.GetChatByUsernames(username, targetUsername)
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "component": "chat",
        }).Errorf("chat not found with error: %v", err.Error())
        c.JSON(http.StatusBadRequest, gin.H{"error": "chat not found"})
        return
    }

    lastMessage, err := ctrl.chatService.GetLastMessageByChatID(chat.ID)
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "component": "chat",
        }).Errorf("Failed to get last message with error: %v", err.Error())
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get last message"})
        return
    }

    curUsr, err := ctrl.userService.GetUserByUsername(username)
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "component": "chat",
        }).Errorf("Failed to get current user with error: %v", err.Error())
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get current user"})
        return
    }
	log.Println(lastMessage.SenderID, curUsr.ID)
    if lastMessage.SenderID != curUsr.ID {
        // task, err := tasks.NewReadMessagesTask(chat.ID, curUsr.ID)
		// log.Println(string(task.Payload()))
        // if err != nil {
        //     logger.Log.WithFields(logrus.Fields{
        //         "component": "chat",
        //     }).Errorf("Failed to create read messages task with error: %v", err.Error())
        //     c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create read messages task"})
        //     return
        // }

        // inf, err = redis.Client.Enqueue(task)
		// if err != nil {
		// 	log.Panicln(err.Error())
		// }
		// log.Println(inf)
		task, err := tasks.NewReadMessagesTask(chat.ID, lastMessage.SenderID)
		if err != nil {
			logger.Log.WithFields(logrus.Fields{
				"component": "chat",
				"service":   "asynq",
			}).Errorf("server could not start messages reader task with error: %v", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}
		info, err := redis.Client.Enqueue(task)
		if err != nil {
			logger.Log.WithFields(logrus.Fields{
				"component": "chat",
				"service":   "asynq",
			}).Errorf("server could not start messages reader with error: %v", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}
		log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
	}

    messages, err := ctrl.chatService.GetMessagesByIDChat(chat.ID)
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "component": "chat",
        }).Errorf("Failed to get messages with error: %v", err.Error())
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get messages"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "chat_members": chat,
        "messages": messages,
    })
}

// ChatController godoc
// @Summary Get all chats for current user
// @Description Route which return all chats for current user
// @Tags chat
// @Produce json
// @Param Authorization header string true "With the Bearer started"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /chats [get]
func (ctrl *ChatController) GetChatsForSpecUser(c *gin.Context) {
	username := c.MustGet("username").(string)
	user, err := ctrl.userService.GetUserByUsername(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	chats, err := ctrl.chatService.GetUserChats(user.ID)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"component": "chat",
		}).Errorf("Failed to get chats for special user: %v with error: %v", username, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get chats for special user"})
		return
	}
	c.JSON(http.StatusOK, chats)
}


type SendMessageInput struct {
	ChatID     uint   `json:"chat_id" binding:"required"`
	ReceiverID uint   `json:"receiver_id" binding:"required"`
	Message    string `json:"message" binding:"required"`
}

// SendMessage godoc
// @Summary Send a message in a chat
// @Description This endpoint allows an authenticated user to send a message in a specific chat.
// @Tags chat
// @Accept  json
// @Produce  json
// @Param Authorization header string true "With the Bearer started"
// @Param SendMessageInput body SendMessageInput true "Message input data"
// @Success 200 {string} string "Message sent successfully"
// @Failure 400 {object} string "Invalid input data"
// @Failure 404 {object} string "User not found"
// @Failure 500 {object} string "Failed to create the message"
// @Router /chats/message [post]
func (ctrl *ChatController) SendMessage(c *gin.Context) {
	var input SendMessageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"component": "chat",
		}).Errorf("Failed to receive data from client: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	username := c.MustGet("username").(string)
	user, err := ctrl.userService.GetUserByUsername(username)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"component": "chat",
			"username":  username,
		}).Errorf("User not found: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	message := models.Message{
		ChatID:    input.ChatID,
		SenderID:  user.ID,
		ReceiverID: input.ReceiverID,
		Content:   input.Message,
		IsRead:    false,
	}

	if err := ctrl.chatService.CreateMessage(&message); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"component": "chat",
			"chat_id":   input.ChatID,
			"sender_id": user.ID,
			"receiver_id": input.ReceiverID,
		}).Errorf("Failed to create message: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create the message"})
		return 
	}

	c.Status(http.StatusOK)
}
