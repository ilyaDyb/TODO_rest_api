package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/ilyaDyb/go_rest_api/config"
	_ "github.com/ilyaDyb/go_rest_api/docs"
	"github.com/ilyaDyb/go_rest_api/middleware"
	"github.com/ilyaDyb/go_rest_api/pereodictasks"
	"github.com/ilyaDyb/go_rest_api/routes"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)


const redisAddr = "localhost:6379"
var Client *asynq.Client

func init() {
	Client = asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
}

// @title           Swagger REST API
// @version         1.0
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.basic  JWT

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {

	config.Connect()
	router := gin.Default()

	// router.Use(cors.Default())
	router.Use(middleware.CORSMiddleware())
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	routes.TestRoute(router)
	routes.AuthRoute(router)
	routes.UserRoute(router)
	routes.AdminRoute(router)
	go func () {
		if err := router.Run(":8080"); err != nil {
			log.Fatalf("could not run sever: %v", err)
		}
	} ()
	if err := pereodictasks.StartPereodicTasks(); err != nil {
		log.Fatalln(err)
	}
	if err := config.StartRedis(); err != nil {
		log.Fatalln(err)
	}
}
