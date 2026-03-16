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

// BugCache handles Redis operations for bug caching

// we create different cache  structure/files for different domain (bug , user) because their storing data or
//acessing db pattern may  be differnet domain to domain

type BugCache struct {
	redis *redis.Client
}

// NewBugCache creates a new BugCache instance
func NewBugCache(redis *redis.Client) *BugCache {
	return &BugCache{redis: redis}
}

// Get retrieves a bug from cache by ID
func (c *BugCache) Get(ctx context.Context, bugID int64) (*models.Bug, error) {
	key := fmt.Sprintf("bug:%d", bugID)
   
	cached, err := c.redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			metrics.RedisCacheMisses.Inc()
		}
		return nil, err // Return error to indicate not found or connection issue
	}

	var bug models.Bug
	if err := json.Unmarshal([]byte(cached), &bug); err != nil {
		return nil, err
	}
     fmt.Println("this is the bug" , bug)
	return &bug, nil
}

// Set stores a bug in cache
func (c *BugCache) Set(ctx context.Context, bugID int64, bug *models.Bug) error {
	key := fmt.Sprintf("bug:%d", bugID)
	ttl := 1 * time.Hour  // Cache for 1 hour instead of 10 seconds
	data, err := json.Marshal(bug)
	if err != nil {
		return err
	}
	return c.redis.Set(ctx, key, data, ttl).Err()
}

// Delete removes a bug from cache
func (c *BugCache) Delete(ctx context.Context, bugID int64) error {
	key := fmt.Sprintf("bug:%d", bugID)
	return c.redis.Del(ctx, key).Err()
}

func (c *BugCache) ClearAll(ctx context.Context) error {
	return c.redis.FlushDB(ctx).Err()
}