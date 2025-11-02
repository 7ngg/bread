package cache

import (
	"github.com/7ngg/bread/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisConnection(cfg *config.RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
		Protocol: cfg.Protocol,
	})

	return client
}
