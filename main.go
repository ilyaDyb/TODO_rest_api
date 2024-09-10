package main

import (
	"log"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/ilyaDyb/go_rest_api/config"
	"github.com/ilyaDyb/go_rest_api/config/redis"
	_ "github.com/ilyaDyb/go_rest_api/docs"
	"github.com/ilyaDyb/go_rest_api/logger"
	"github.com/ilyaDyb/go_rest_api/middleware"
	"github.com/ilyaDyb/go_rest_api/pereodictasks"
	"github.com/ilyaDyb/go_rest_api/routes"
	"github.com/ilyaDyb/go_rest_api/ws"
	"github.com/sirupsen/logrus"
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
	esConfig := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	}
	client, err := elasticsearch.NewClient(esConfig)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"service": "elasticsearch",
		}).Fatalf("could not run elasticsearch sever: %v", err)
		panic(err)
	}
	logger.InitLogger(client)
	
	config.Connect()
	router := gin.Default()
	
	// router.Use(cors.Default())
	router.Use(middleware.CORSMiddleware())
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	routes.TestRoute(router)
	routes.AuthRoute(router)
	routes.UserRoute(router)
	routes.ChatRoute(router)
	routes.AdminRoute(router)

	ws.RegisterWsRoutes(router)
	go ws.HubInstance.Run()

	log.Println("Calling run server...")
	go func() {
		if err := router.Run(":8080"); err != nil {
			logger.Log.WithFields(logrus.Fields{
				"service": "server",
			}).Fatalf("could not run sever: %v", err)
			log.Fatalf("could not run sever: %v", err)
		}
	}()
	
	log.Println("Calling StartPereodicTasks...")
	if err := pereodictasks.StartPereodicTasks(); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"service": "asynq",
		}).Fatalf("could not run asynq sever: %v", err)
		log.Fatalln(err)
	}
	log.Println("Calling StartRedis...")
	if err := redis.StartRedis(); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"service": "redis",
		}).Fatalf("could not run redis sever: %v", err)
		log.Fatalln(err)
	}
	logger.Log.Info("All applications was started")
}
