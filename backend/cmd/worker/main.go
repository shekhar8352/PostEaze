package main

import (
	"log"

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

	// Initialize Asynq Server
	if err := tasks.InitServer(); err != nil {
		log.Fatalf("Failed to init asynq server: %v", err)
	}

	// Initialize Asynq Scheduler (for cron jobs)
	if err := tasks.InitScheduler(); err != nil {
		log.Fatalf("Failed to init asynq scheduler: %v", err)
	}

	// Start Scheduler in a goroutine
	go tasks.StartScheduler()

	log.Println("Starting Asynq Worker...")
	// Start Server (Blocking)
	tasks.StartServer()
}
