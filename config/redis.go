package config
import (
	"log"

	"github.com/hibiken/asynq"
	"github.com/ilyaDyb/go_rest_api/tasks"

	"github.com/redis/go-redis/v9"
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
	mux.HandleFunc("email:deliver", tasks.HandleEmailDeliveryTask)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run Asynq server: %v", err)
		return err
	}
	return nil
}