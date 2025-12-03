package main

import (
	"context"
	"log"

	"github.com/shekhar8352/PostEaze/constants"
	"github.com/shekhar8352/PostEaze/tasks"
	"github.com/shekhar8352/PostEaze/utils/configs"
	"github.com/shekhar8352/PostEaze/utils/database"
	"github.com/shekhar8352/PostEaze/utils/encryption"
	"github.com/shekhar8352/PostEaze/utils/env"
	"github.com/shekhar8352/PostEaze/utils/flags"
)

func main() {
	ctx := context.Background()
	env.InitEnv()

	// Initialize configs
	configNames := []string{constants.DatabaseConfig}
	err := configs.InitDev(flags.BaseConfigPath(), configNames...)
	if err != nil {
		log.Fatalf("Failed to init configs: %v", err)
	}

	// Initialize database connection
	driverName, err := configs.Get().GetString(constants.DatabaseConfig, constants.DatabaseDriverNameConfigKey)
	if err != nil {
		log.Fatalf("Failed to get database driver name: %v", err)
	}
	urlString, err := configs.Get().GetString(constants.DatabaseConfig, constants.DatabaseURLConfigKey)
	if err != nil {
		log.Fatalf("Failed to get database url: %v", err)
	}
	// Apply environment variable substitution
	url := env.ApplyEnvironmentToString(urlString)

	dbConfig := database.Config{
		DriverName: driverName,
		URL:        url,
	}
	if err := database.Init(ctx, dbConfig); err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}
	log.Println("Database initialized successfully")

	// Initialize encryption
	if err := encryption.Init(); err != nil {
		log.Fatalf("Failed to init encryption: %v", err)
	}
	log.Println("Encryption initialized successfully")

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
