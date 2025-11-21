package tasks

import (
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/shekhar8352/PostEaze/constants"
	"github.com/shekhar8352/PostEaze/utils/configs"
	"github.com/shekhar8352/PostEaze/utils/env"
)

var client *asynq.Client

// InitClient initializes the Asynq client
func InitClient() error {
	host, _ := configs.Get().GetString(constants.DatabaseConfig, constants.RedisHostConfigKey)
	port, _ := configs.Get().GetString(constants.DatabaseConfig, constants.RedisPortConfigKey)
	password, _ := configs.Get().GetString(constants.DatabaseConfig, constants.RedisPasswordConfigKey)
	db := configs.Get().GetIntD(constants.DatabaseConfig, constants.RedisDBConfigKey, 0)

	addr := fmt.Sprintf("%s:%s", host, port)
	addr = env.ApplyEnvironmentToString(addr)
	password = env.ApplyEnvironmentToString(password)

	// Default to localhost:6379 if addr is empty/invalid
	if addr == ":" || addr == ":0" {
		addr = "localhost:6379"
	}

	client = asynq.NewClient(asynq.RedisClientOpt{
		Addr:     addr,
		Password: password,
		DB:       int(db),
	})

	return nil
}

// GetClient returns the Asynq client instance
func GetClient() *asynq.Client {
	return client
}

// EnqueueTask enqueues a task to a specific queue
func EnqueueTask(task *asynq.Task, queue string, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	opts = append(opts, asynq.Queue(queue))
	return client.Enqueue(task, opts...)
}
