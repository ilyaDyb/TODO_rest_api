package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ilyaDyb/go_rest_api/controller"
	"github.com/ilyaDyb/go_rest_api/middleware"
)

func FeedRoute(router *gin.Engine) {
	authorized := router.Group("/feed")
	authorized.Use(middleware.JWTAuthMiddleware())
	{
		authorized.POST("/grade", controller.GradeProfile)
		authorized.GET("/get-profiles", controller.GetProfiles)
	}
}