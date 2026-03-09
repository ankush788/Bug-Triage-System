package pkg

import (
	"context"
	"sync"

	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// RateLimiter provides in-memory rate limiting using token bucket algorithm
// Each identifier gets its own limiter with configurable rate and burst.
// Uses golang.org/x/time/rate for standard Go rate limiting.

type RateLimiter struct {
	limiters sync.Map // map[string]*rate.Limiter
	logger   *zap.Logger
	rate     rate.Limit // requests per second
	burst    int        // burst capacity
}

func NewRateLimiter(logger *zap.Logger) *RateLimiter {
	return &RateLimiter{
		logger: logger,
		rate:   1 , // refilling rate 10 request 60 second
		burst:  10, // maximum token
	}
}


// AllowRequest checks if a request from the given identifier is allowed
// Returns true if within rate limit, false otherwise
func (rl *RateLimiter) AllowRequest(ctx context.Context, identifier string) bool {
	limiter, _ := rl.limiters.LoadOrStore(identifier, rate.NewLimiter(rl.rate, rl.burst))
	allowed := limiter.(*rate.Limiter).Allow()
	if !allowed {
		rl.logger.Info("rate limit exceeded", zap.String("identifier", identifier))
	}
	return allowed
}

// RemainingTokens returns an approximate number of tokens remaining for an identifier
func (rl *RateLimiter) RemainingTokens(ctx context.Context, identifier string) (int64, error) {
	// Note: golang.org/x/time/rate doesn't expose internal token count
	// Return burst capacity as approximation
	return int64(rl.burst), nil
}

// Reset clears the rate limit for an identifier
func (rl *RateLimiter) Reset(ctx context.Context, identifier string) error {
	rl.limiters.Delete(identifier)
	return nil
}
