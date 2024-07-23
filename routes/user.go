package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ilyaDyb/go_rest_api/controller"
	"github.com/ilyaDyb/go_rest_api/middleware"
)

func TestRoute(router *gin.Engine) {
	router.GET("/test/:name", controller.UserNameController)
	router.GET("/test/:name/*action", controller.UserNameActionController)
	router.GET("/welcome", controller.WelcomeController)
	router.POST("/form_post", controller.FormPostController)
	router.POST("/queryform_post", controller.QueryFormPostController)
	router.MaxMultipartMemory = 8 << 20
	router.POST("/upload", controller.UploadFile)
	router.GET("/testing", controller.TestingValidate)
	router.GET("/test/queries", controller.TestQueries)
}
func UserRoute(router *gin.Engine) {
	authorized := router.Group("/u")
	authorized.Use(middleware.JWTAuthMiddleware())
	{
		authorized.GET("/profile/*username", controller.ProfileController)
		authorized.PUT("/profile", controller.EditProfileController)
	}
}