package redis

import (
	"context"
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/shekhar8352/PostEaze/constants"
	"github.com/shekhar8352/PostEaze/utils/configs"
	"github.com/shekhar8352/PostEaze/utils/env"
)

var client *redis.Client

// Init initializes the Redis client
func Init(ctx context.Context) error {
	host, err := configs.Get().GetString(constants.DatabaseConfig, constants.RedisHostConfigKey)
	if err != nil {
		return fmt.Errorf("failed to get redis host: %w", err)
	}

	port, err := configs.Get().GetString(constants.DatabaseConfig, constants.RedisPortConfigKey)
	if err != nil {
		return fmt.Errorf("failed to get redis port: %w", err)
	}

	password, _ := configs.Get().GetString(constants.DatabaseConfig, constants.RedisPasswordConfigKey)
	// Password might be empty, so we ignore the error or check if it's required based on env

	db := configs.Get().GetIntD(constants.DatabaseConfig, constants.RedisDBConfigKey, 0)

	addr := fmt.Sprintf("%s:%s", host, port)
	// Apply environment variable substitution if needed (e.g. if host is "${REDIS_HOST}")
	addr = env.ApplyEnvironmentToString(addr)
	password = env.ApplyEnvironmentToString(password)

	// Default to localhost:6379 if addr is just ":" (meaning host/port were empty)
	if addr == ":" || addr == ":0" {
		addr = "localhost:6379"
	} else if strings.HasPrefix(addr, ":") {
		addr = "localhost" + addr
	} else if strings.HasSuffix(addr, ":") {
		addr = addr + "6379"
	}

	client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       int(db),
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	return nil
}

// GetClient returns the Redis client
func GetClient() *redis.Client {
	return client
}

// Ping checks connectivity to Redis (PING command).
func Ping(ctx context.Context) error {
	if client == nil {
		return fmt.Errorf("redis not initialized")
	}
	return client.Ping(ctx).Err()
}
