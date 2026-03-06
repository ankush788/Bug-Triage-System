package pkg

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// RateLimiter uses Redis for distributed rate limiting
// Implements token-bucket algorithm
// Each identifier has a bucket with a fixed capacity (maxTokens) and a refill rate.
// Tokens accumulate at the refill rate (per second) up to the capacity.
// Each request consumes one token. If tokens are available, the request is allowed.
// The entire decision is performed atomically by a Lua script.
//
// This design allows controlled burst traffic while maintaining average rate limits.

type RateLimiter struct {
    client          *redis.Client
    logger          *zap.Logger
    maxTokens       int64 // bucket capacity
    refillRate      int64 // tokens per second that are added
    requestCost     int64 // tokens consumed per request
}

func NewRateLimiter(client *redis.Client, logger *zap.Logger) *RateLimiter {
    return &RateLimiter{
        client:      client,
        logger:      logger,
        maxTokens:   100,  // maximum tokens in bucket (burst capacity)
        refillRate:  10,   // tokens added per second
        requestCost: 1,    // tokens consumed per request
    }
}

// AllowRequest checks if a request from the given identifier is allowed
// Returns true if within rate limit, false otherwise
func (rl *RateLimiter) AllowRequest(ctx context.Context, identifier string) bool {
	key := fmt.Sprintf("rate_limit:%s", identifier)

    // Lua script implements token bucket atomically
    // ARGV[1] = capacity, ARGV[2] = refillRate, ARGV[3] = now (milliseconds), ARGV[4] = requestCost
    // it stores a hash {tokens=<current>, last=<timestamp>} and updates it
    script := redis.NewScript(`
        local key = KEYS[1]
        local capacity = tonumber(ARGV[1])
        local refill_rate = tonumber(ARGV[2])
        local now = tonumber(ARGV[3])
        local request_cost = tonumber(ARGV[4])

        local bucket = redis.call('hmget', key, 'tokens', 'last')
        local tokens = tonumber(bucket[1]) or capacity
        local last = tonumber(bucket[2]) or now

        -- calculate tokens added since last check (in milliseconds)
        local elapsed = now - last
        local refilled = (elapsed / 1000) * refill_rate
        tokens = tokens + refilled
        if tokens > capacity then
            tokens = capacity
        end

        if tokens >= request_cost then
            -- request allowed, consume tokens
            tokens = tokens - request_cost
            redis.call('hmset', key, 'tokens', tokens, 'last', now)
            redis.call('expire', key, 3600)
            return 1
        else
            -- request denied, update timestamp for next calculation
            redis.call('hmset', key, 'tokens', tokens, 'last', now)
            redis.call('expire', key, 3600)
            return 0
        end
    `)

    now := time.Now().UnixMilli()
    result, err := script.Run(ctx, rl.client, []string{key}, rl.maxTokens, rl.refillRate, now, rl.requestCost).Int64()
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

// RemainingTokens returns an approximate number of tokens remaining for an identifier
func (rl *RateLimiter) RemainingTokens(ctx context.Context, identifier string) (int64, error) {
    key := fmt.Sprintf("rate_limit:%s", identifier)

    vals, err := rl.client.HMGet(ctx, key, "tokens", "last").Result()
    if err != nil {
        if err == redis.Nil {
            return rl.maxTokens, nil
        }
        return 0, err
    }

    var tokens float64 = float64(rl.maxTokens)
    var last int64 = time.Now().UnixMilli()
    
    if vals[0] != nil {
        if str, ok := vals[0].(string); ok {
            if val, err := strconv.ParseFloat(str, 64); err == nil {
                tokens = val
            }
        }
    }
    if vals[1] != nil {
        if str, ok := vals[1].(string); ok {
            if val, err := strconv.ParseInt(str, 10, 64); err == nil {
                last = val
            }
        }
    }

    // refill since last check (in milliseconds)
    elapsed := time.Now().UnixMilli() - last
    refilled := (float64(elapsed) / 1000.0) * float64(rl.refillRate)
    tokens += refilled
    if tokens > float64(rl.maxTokens) {
        tokens = float64(rl.maxTokens)
    }
    if tokens < 0 {
        tokens = 0
    }

    return int64(tokens), nil
}

// Reset clears the rate limit for an identifier
func (rl *RateLimiter) Reset(ctx context.Context, identifier string) error {
	key := fmt.Sprintf("rate_limit:%s", identifier)
	return rl.client.Del(ctx, key).Err()
}
