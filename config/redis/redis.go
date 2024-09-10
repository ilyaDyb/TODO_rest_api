package redis

import (
	"log"

	"github.com/hibiken/asynq"
	"github.com/ilyaDyb/go_rest_api/logger"
	"github.com/ilyaDyb/go_rest_api/tasks"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

const RedisAddr = "localhost:6379"

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
	mux.HandleFunc("email:deliver", tasks.HandleEmailDeliveryTask)
	mux.HandleFunc("messages:reader", tasks.HandleReadMessagesTask)
	
	log.Println("Starting Asynq server...")
	if err := srv.Run(mux); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"service": "async",
		}).Fatalf("async start failed with error: %v", err.Error())
		log.Fatalf("could not run Asynq server: %v", err)
		return err
	}
	return nil
}

