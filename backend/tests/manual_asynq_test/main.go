package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
	"github.com/shekhar8352/PostEaze/constants"
	"github.com/shekhar8352/PostEaze/tasks"
	"github.com/shekhar8352/PostEaze/utils/configs"
	"github.com/shekhar8352/PostEaze/utils/env"
	"github.com/shekhar8352/PostEaze/utils/flags"
)

func main() {
	// ctx := context.Background()
	env.InitEnv()

	// Initialize configs
	configNames := []string{constants.DatabaseConfig}
	err := configs.InitDev(flags.BaseConfigPath(), configNames...)
	if err != nil {
		log.Fatalf("Failed to init configs: %v", err)
	}

	// Initialize Asynq Client
	err = tasks.InitClient()
	if err != nil {
		log.Fatalf("Failed to init asynq client: %v", err)
	}

	// Create a task
	payload, _ := json.Marshal(tasks.LogMessagePayload{Message: "Hello from Manual Test!"})
	task := asynq.NewTask(tasks.TypeLogMessage, payload)

	// Enqueue task to 'fast' queue
	info, err := tasks.EnqueueTask(task, tasks.QueueFast)
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	fmt.Printf("Enqueued task: id=%s queue=%s\n", info.ID, info.Queue)

	// Note: To see it processed, the main backend server needs to be running,
	// OR we can start a temporary server here just for verification.
	// For this test, we'll just verify enqueueing works.
	// The user can verify processing by running the main app.
	fmt.Println("Task enqueued successfully. Run the main application to see it processed.")
}
