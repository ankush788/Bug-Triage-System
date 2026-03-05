package pkg

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// RateLimiter uses Redis for distributed rate limiting
// Implements token bucket algorithm
type RateLimiter struct {
	client      *redis.Client
	logger      *zap.Logger
	maxTokens   int64
	refillRate  int64         // tokens per second
	refillEvery time.Duration // how often to check for refill
}

func NewRateLimiter(client *redis.Client, logger *zap.Logger) *RateLimiter {
	return &RateLimiter{
		client:      client,
		logger:      logger,
		maxTokens:   100,
		refillRate:  10,           // 10 tokens per second
		refillEvery: time.Second,
	}
}

// AllowRequest checks if a request from the given identifier is allowed
// Returns true if within rate limit, false otherwise
func (rl *RateLimiter) AllowRequest(ctx context.Context, identifier string) bool {
	key := fmt.Sprintf("rate_limit:%s", identifier)

	// Lua script: check and decrement tokens atomically
	script := redis.NewScript(`
        local key = KEYS[1]
        local limit = tonumber(ARGV[1])
        local current = redis.call('get', key)
        
        if not current then
            redis.call('set', key, limit - 1)
            redis.call('expire', key, 60)
            return 1
        end
        
        current = tonumber(current)
        if current > 0 then
            redis.call('decr', key)
            return 1
        end
        
        return 0
    `)

	result, err := script.Run(ctx, rl.client, []string{key}, rl.maxTokens).Int64()
	if err != nil {
		rl.logger.Error("rate limiter redis error", zap.Error(err))
		// In case of Redis failure, allow request (fail open)
		return true
	}

	allowed := result == 1

	if !allowed {
		rl.logger.Debug("rate limit exceeded", zap.String("identifier", identifier))
	}

	return allowed
}

// RemainingTokens returns the number of tokens remaining for an identifier
func (rl *RateLimiter) RemainingTokens(ctx context.Context, identifier string) (int64, error) {
	key := fmt.Sprintf("rate_limit:%s", identifier)

	val, err := rl.client.Get(ctx, key).Int64()
	if err == redis.Nil {
		return rl.maxTokens, nil
	}
	if err != nil {
		return 0, err
	}

	return val, nil
}

// Reset clears the rate limit for an identifier
func (rl *RateLimiter) Reset(ctx context.Context, identifier string) error {
	key := fmt.Sprintf("rate_limit:%s", identifier)
	return rl.client.Del(ctx, key).Err()
}
