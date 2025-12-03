package tasks

import (
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
	"github.com/shekhar8352/PostEaze/constants"
	"github.com/shekhar8352/PostEaze/utils/configs"
	"github.com/shekhar8352/PostEaze/utils/env"
)

var scheduler *asynq.Scheduler

// InitScheduler initializes the Asynq scheduler
func InitScheduler() error {
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

	scheduler = asynq.NewScheduler(
		asynq.RedisClientOpt{
			Addr:     addr,
			Password: password,
			DB:       int(db),
		},
		&asynq.SchedulerOpts{
			Location: time.UTC,
		},
	)

	// Register Instagram profile sync job to run every 6 hours (using slow queue for background jobs)
	syncTask := asynq.NewTask(TypeSyncInstagramProfiles, nil)
	if _, err := scheduler.Register("@every 6h", syncTask, asynq.Queue(QueueSlow)); err != nil {
		log.Printf("Warning: Failed to register Instagram profile sync job: %v", err)
	} else {
		log.Println("Registered Instagram profile sync job to run every 6 hours")
	}

	return nil
}

// StartScheduler starts the scheduler
func StartScheduler() {
	if err := scheduler.Run(); err != nil {
		log.Fatalf("could not run scheduler: %v", err)
	}
}
