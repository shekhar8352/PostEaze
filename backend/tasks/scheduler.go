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

// Instagram periodic sync specs (asynq "@every" duration). Sub-hour cadence; tune if you hit Graph API rate limits.
const (
	instagramProfileSyncSpec   = "@every 1h"
	instagramPostsSyncSpec     = "@every 30m"
	instagramAnalyticsSyncSpec = "@every 30m"
)

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

	// Instagram sync jobs (slow queue). Restart scheduler after changing constants above.
	syncTask := asynq.NewTask(TypeSyncInstagramProfiles, nil)
	if _, err := scheduler.Register(instagramProfileSyncSpec, syncTask, asynq.Queue(QueueSlow)); err != nil {
		log.Printf("Warning: Failed to register Instagram profile sync job: %v", err)
	} else {
		log.Printf("Registered Instagram profile sync job (%s)", instagramProfileSyncSpec)
	}

	postsTask := asynq.NewTask(TypeSyncInstagramPosts, nil)
	if _, err := scheduler.Register(instagramPostsSyncSpec, postsTask, asynq.Queue(QueueSlow)); err != nil {
		log.Printf("Warning: Failed to register Instagram posts sync job: %v", err)
	} else {
		log.Printf("Registered Instagram posts sync job (%s)", instagramPostsSyncSpec)
	}

	analyticsTask := asynq.NewTask(TypeSyncInstagramAnalytics, nil)
	if _, err := scheduler.Register(instagramAnalyticsSyncSpec, analyticsTask, asynq.Queue(QueueSlow)); err != nil {
		log.Printf("Warning: Failed to register Instagram analytics sync job: %v", err)
	} else {
		log.Printf("Registered Instagram analytics sync job (%s)", instagramAnalyticsSyncSpec)
	}

	// Register period snapshot task to run daily
	snapshotTask := asynq.NewTask(TypePeriodSnapshot, nil)
	if _, err := scheduler.Register("@daily", snapshotTask, asynq.Queue(QueueSlow)); err != nil {
		log.Printf("Warning: Failed to register period snapshot job: %v", err)
	} else {
		log.Println("Registered period snapshot job to run daily")
	}

	return nil
}

// StartScheduler starts the scheduler
func StartScheduler() {
	if err := scheduler.Run(); err != nil {
		log.Fatalf("could not run scheduler: %v", err)
	}
}
