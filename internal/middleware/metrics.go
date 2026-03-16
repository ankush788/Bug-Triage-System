package middleware

import (
	"bug_triage/internal/metrics"
	"time"

	"github.com/gin-gonic/gin"
)

// MetricsMiddleware collects request latency metrics.
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		endpoint := c.FullPath()
		if endpoint == "" {
			endpoint = c.Request.URL.Path
		}

		duration := time.Since(start).Seconds()
		metrics.HTTPLatency.WithLabelValues(
			c.Request.Method,
			endpoint,
		).Observe(duration)
	}
}
