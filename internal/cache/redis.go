package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// NewRedisClient initializes and returns a Redis client
func NewRedisClient(addr, password string, RedisDB int,  log *zap.Logger) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB : RedisDB,
	})

	// Test connection
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Error("redis ping failed", zap.Error(err))
		return nil, err
	}

	log.Info("redis connected")
	return redisClient, nil
}
