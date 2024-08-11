package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ilyaDyb/go_rest_api/config"
	"github.com/ilyaDyb/go_rest_api/controller"
	"github.com/ilyaDyb/go_rest_api/repository"
	"github.com/ilyaDyb/go_rest_api/service"
)

func AdminRoute(route *gin.Engine) {
	adminGroup := route.Group("/admin")
	db := config.DB

	adminRepo := repository.NewPostgresUserRepo(db)
	chatRepo := repository.NewPostgresChatRepo(db)

	adminService := service.NewUserService(adminRepo)
	chatService := service.NewChatService(chatRepo)

	adminController := controller.NewAdminController(adminService, chatService)
	{
		adminGroup.GET("/users", adminController.UsersList)
		adminGroup.GET("/user/:id", adminController.GetUser)
		adminGroup.POST("/user", adminController.CreateUser)
		adminGroup.PUT("/user/:id", adminController.UpdateUser)
		adminGroup.DELETE("/user/:id", adminController.DeleteUser)
		
		// adminGroup
		adminGroup.GET("/chats", adminController.GetAllChats)
	}
}