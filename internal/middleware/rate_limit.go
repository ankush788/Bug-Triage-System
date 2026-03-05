package middleware

import (
	"net/http"

	"bug_triage/internal/pkg"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RateLimitMiddleware applies rate limiting to requests
// Uses authenticated user ID if available, otherwise uses IP address
func RateLimitMiddleware(limiter *pkg.RateLimiter, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		identifier := c.ClientIP()

		// Use user ID as identifier if authenticated
		if userID, ok := GetUserID(c); ok {
			identifier = userIDToString(userID)
		}

		if !limiter.AllowRequest(c.Request.Context(), identifier) {
			logger.Warn("rate limit exceeded", zap.String("identifier", identifier))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func userIDToString(userID int64) string {
	// Simple conversion - in production might use a more sophisticated method
	return "user:" + string(rune(userID))
}
