package tasks

import (
	"fmt"
	"log"

	"github.com/hibiken/asynq"
	"github.com/shekhar8352/PostEaze/constants"
	"github.com/shekhar8352/PostEaze/utils/configs"
	"github.com/shekhar8352/PostEaze/utils/env"
)

var server *asynq.Server
var mux *asynq.ServeMux

// InitServer initializes the Asynq server with prioritized queues
func InitServer() error {
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

	server = asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     addr,
			Password: password,
			DB:       int(db),
		},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				QueueFast:   6,
				QueueMedium: 3,
				QueueSlow:   1,
			},
		},
	)

	mux = asynq.NewServeMux()
	RegisterHandlers(mux)

	return nil
}

// StartServer starts the Asynq server
func StartServer() {
	if err := server.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
