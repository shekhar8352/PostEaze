package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/shekhar8352/PostEaze/constants"
	"github.com/shekhar8352/PostEaze/services"
	"github.com/shekhar8352/PostEaze/utils/configs"
	"github.com/shekhar8352/PostEaze/utils/env"
	"github.com/shekhar8352/PostEaze/utils/flags"
	"github.com/shekhar8352/PostEaze/utils/redis"
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

	// Initialize Redis Util
	err = redis.Init(ctx)
	if err != nil {
		log.Fatalf("Failed to init redis util: %v", err)
	}

	// Initialize Service
	redisService := services.NewRedisService()

	// Test 1: Set with default timeout (no argument)
	key1 := "test_service_default_timeout"
	val1 := "value1"
	fmt.Printf("Testing Set with default timeout for key: %s\n", key1)
	err = redisService.Set(ctx, key1, val1)
	if err != nil {
		log.Fatalf("Failed to set key1: %v", err)
	}

	// Verify TTL is around 24 hours
	ttl, err := redis.GetClient().TTL(ctx, key1).Result()
	if err != nil {
		log.Fatalf("Failed to get TTL for key1: %v", err)
	}
	fmt.Printf("TTL for key1: %v (Expected ~24h)\n", ttl)
	if ttl < 23*time.Hour || ttl > 25*time.Hour {
		log.Fatalf("TTL check failed! Expected ~24h, got %v", ttl)
	}

	// Test 2: Get
	fmt.Printf("Testing Get for key: %s\n", key1)
	gotVal, err := redisService.Get(ctx, key1)
	if err != nil {
		log.Fatalf("Failed to get key1: %v", err)
	}
	if gotVal != val1 {
		log.Fatalf("Value mismatch! Expected %s, got %s", val1, gotVal)
	}
	fmt.Printf("Get successful: %s\n", gotVal)

	// Test 3: Delete
	fmt.Printf("Testing Delete for key: %s\n", key1)
	err = redisService.Delete(ctx, key1)
	if err != nil {
		log.Fatalf("Failed to delete key1: %v", err)
	}

	// Verify deletion
	_, err = redisService.Get(ctx, key1)
	if err == nil {
		log.Fatalf("Key1 should have been deleted but was found!")
	}
	fmt.Printf("Delete successful (key not found as expected)\n")

	fmt.Println("All Redis Service tests passed!")
}
