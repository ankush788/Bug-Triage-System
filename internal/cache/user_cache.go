package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"bug_triage/internal/metrics"
	"bug_triage/internal/models"

	"github.com/redis/go-redis/v9"
)

// UserCache handles Redis operations for user caching
type UserCache struct {
	redis *redis.Client
}

// NewUserCache creates a new UserCache instance
func NewUserCache(redis *redis.Client) *UserCache {
	return &UserCache{redis: redis}
}

// Get retrieves a user from cache by email
func (c *UserCache) Get(ctx context.Context, email string) (*models.User, error) {
	key := fmt.Sprintf("user:%s", email)

	cached, err := c.redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			metrics.RedisCacheMisses.Inc()
		}
		return nil, err // Return error to indicate not found or connection issue
	}

	var user models.User
	if err := json.Unmarshal([]byte(cached), &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// Set stores a user in cache
func (c *UserCache) Set(ctx context.Context, email string, user *models.User) error {
	key := fmt.Sprintf("user:%s", email)
	ttl := 1 * time.Second // Cache for 1 hour
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return c.redis.Set(ctx, key, data, ttl).Err()
}

// Delete removes a user from cache
func (c *UserCache) Delete(ctx context.Context, email string) error {
	key := fmt.Sprintf("user:%s", email)
	return c.redis.Del(ctx, key).Err()
}
