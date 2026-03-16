package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// HTTP metrics
var (
	HTTPLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Request latency",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
)

// Database metrics
var (
	DBQueryDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Database query duration",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation", "table"},
	)
)

// Redis metrics
var (
	RedisCacheMisses = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "redis_cache_misses_total",
			Help: "Redis cache misses",
		},
	)
)

// AI / RAG metrics
var (
	AILatency = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "ai_analysis_duration_seconds",
			Help:    "AI analysis duration",
			Buckets: prometheus.DefBuckets,
		},
	)
)

// Init registers all metrics with the global Prometheus registry.
func Init() {
	prometheus.MustRegister(
		HTTPLatency,
		RedisCacheMisses,
		AILatency,
		DBQueryDuration,
	)
}
