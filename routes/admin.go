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
	adminService := service.NewUserService(adminRepo)
	adminController := controller.NewAdminController(adminService)
	{
		adminGroup.GET("/users", adminController.UsersList)
		adminGroup.GET("/user/:id", adminController.GetPutPostDeleteUser)
		adminGroup.POST("/user", adminController.GetPutPostDeleteUser)
		adminGroup.PUT("/user/:id", adminController.GetPutPostDeleteUser)
		adminGroup.DELETE("/user/:id", adminController.GetPutPostDeleteUser)
	}
}