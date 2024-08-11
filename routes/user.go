package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ilyaDyb/go_rest_api/config"
	"github.com/ilyaDyb/go_rest_api/controller"
	"github.com/ilyaDyb/go_rest_api/middleware"
	"github.com/ilyaDyb/go_rest_api/repository"
	"github.com/ilyaDyb/go_rest_api/service"
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
	db := config.DB

	userRepo := repository.NewPostgresUserRepo(db)
	chatRepo := repository.NewPostgresChatRepo(db)

	userService := service.NewUserService(userRepo)
	chatService := service.NewChatService(chatRepo)

	userController := controller.NewUserController(userService, chatService)

	authorized := router.Group("/u")
	authorized.Use(middleware.JWTAuthMiddleware())
	{
		authorized.GET("/profile/*username", userController.ProfileController)
		authorized.PUT("/profile", userController.EditProfileController)
		authorized.PATCH("/set-as-preview/:photo_id", userController.SetAsPriviewController)
		authorized.PATCH("/save-location", userController.SaveLocationController)
		authorized.PATCH("/set-coordinates", userController.SetCoordinatesController)
		authorized.GET("/liked-by-users", userController.LikedByUsersController)
		authorized.POST("/grade", userController.GradeProfileController)
		authorized.GET("/get-profiles", userController.GetProfilesController)
	}
}