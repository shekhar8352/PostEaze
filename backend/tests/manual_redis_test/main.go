package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/shekhar8352/PostEaze/constants"
	"github.com/shekhar8352/PostEaze/utils/configs"
	"github.com/shekhar8352/PostEaze/utils/env"
	"github.com/shekhar8352/PostEaze/utils/flags"
	"github.com/shekhar8352/PostEaze/utils/redis"
)

func main() {
	ctx := context.Background()
	env.InitEnv()

	// Initialize configs (assuming dev mode for test)
	configNames := []string{constants.DatabaseConfig}
	err := configs.InitDev(flags.BaseConfigPath(), configNames...)
	if err != nil {
		log.Fatalf("Failed to init configs: %v", err)
	}

	// Initialize Redis
	err = redis.Init(ctx)
	if err != nil {
		log.Fatalf("Failed to init redis: %v", err)
	}

	client := redis.GetClient()

	// Set a value
	err = client.Set(ctx, "test_key", "hello redis", 10*time.Second).Err()
	if err != nil {
		log.Fatalf("Failed to set value: %v", err)
	}

	// Get the value
	val, err := client.Get(ctx, "test_key").Result()
	if err != nil {
		log.Fatalf("Failed to get value: %v", err)
	}

	fmt.Printf("Redis test successful! Got value: %s\n", val)
}
