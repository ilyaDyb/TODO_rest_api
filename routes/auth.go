package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ilyaDyb/go_rest_api/controller"
)

func AuthRoute(router *gin.Engine)  {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/registration", controller.Register)
		authGroup.POST("/login", controller.Login)
	}
}