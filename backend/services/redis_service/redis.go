package redis_service

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	redisUtil "github.com/shekhar8352/PostEaze/utils/redis"
)

type RedisService interface {
	Set(ctx context.Context, key string, value interface{}, expiration ...time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

type RedisClientService struct {
	client *redis.Client
}

func NewRedisService() *RedisClientService {
	return &RedisClientService{
		client: redisUtil.GetClient(),
	}
}

func (s *RedisClientService) Set(ctx context.Context, key string, value interface{}, expiration ...time.Duration) error {
	exp := 24 * time.Hour
	if len(expiration) > 0 {
		exp = expiration[0]
	}
	return s.client.Set(ctx, key, value, exp).Err()
}

func (s *RedisClientService) Get(ctx context.Context, key string) (string, error) {
	return s.client.Get(ctx, key).Result()
}

func (s *RedisClientService) Delete(ctx context.Context, key string) error {
	return s.client.Del(ctx, key).Err()
}
