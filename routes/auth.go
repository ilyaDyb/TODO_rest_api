package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ilyaDyb/go_rest_api/config"
	"github.com/ilyaDyb/go_rest_api/controller"
	"github.com/ilyaDyb/go_rest_api/repository"
	"github.com/ilyaDyb/go_rest_api/service"
)

func AuthRoute(router *gin.Engine) {
	authGroup := router.Group("/auth")
	db := config.DB
	authRepo := repository.NewPostgresUserRepo(db)
	authService := service.NewUserService(authRepo)
	authController := controller.NewAuthController(authService)
	{
		authGroup.POST("/registration", authController.RegistrationController)
		authGroup.POST("/login", authController.LoginController)
		authGroup.GET("/confirm", authController.ConfirmEmailController)
		authGroup.POST("/refresh", authController.RefreshController)
		authGroup.POST("/drop-password", authController.DropPasswordController)
		authGroup.POST("/change-password", authController.ChangePassword)
	}
}