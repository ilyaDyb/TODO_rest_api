package config

import (
	"log"

	"github.com/hibiken/asynq"
	"github.com/ilyaDyb/go_rest_api/utils"
	"github.com/redis/go-redis/v9"
)

const (
	DefaultUploadPath = "./uploads/"
	UserPhotoPath     = DefaultUploadPath + "user_photos/"
	RedisAddr         = "localhost:6379"
	ServerHost		  = "localhost:8080"
	ServerProtocol	  = "http://"
)

var (
	Client *asynq.Client
	RedisClient *redis.Client
)

func init() {
	Client = asynq.NewClient(asynq.RedisClientOpt{Addr: RedisAddr})
	RedisClient = redis.NewClient(&redis.Options{Addr: RedisAddr})
}

func StartRedis() error {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: RedisAddr},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"default":  6,
				"critical": 3,
				"low":      1,
			},
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc("email:deliver", utils.HandleEmailDeliveryTask)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run Asynq server: %v", err)
		return err
	}
	return nil
}