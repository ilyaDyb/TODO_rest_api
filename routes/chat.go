package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ilyaDyb/go_rest_api/config"
	"github.com/ilyaDyb/go_rest_api/controller"
	"github.com/ilyaDyb/go_rest_api/middleware"
	"github.com/ilyaDyb/go_rest_api/repository"
	"github.com/ilyaDyb/go_rest_api/service"
)

func ChatRoute(router *gin.Engine) {
	db := config.DB

	userRepo := repository.NewPostgresUserRepo(db)
	chatRepo := repository.NewPostgresChatRepo(db)

	userService := service.NewUserService(userRepo)
	chatService := service.NewChatService(chatRepo)
	
	chatController := controller.NewChatController(chatService, userService)

	chatGroup := router.Group("/chats")
	chatGroup.Use(middleware.JWTAuthMiddleware())
	{
		chatGroup.GET("", chatController.GetChatsForSpecUser)
		chatGroup.GET("/:username", chatController.ChatController)
		chatGroup.POST("/message", chatController.SendMessage)
	}
}
