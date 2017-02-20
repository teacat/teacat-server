package middleware

import (
	"../../module/metrics"
	"github.com/gin-gonic/gin"
)

// Metrics is a middleware function that initializes the metrics and attaches to
// the context of every request context.
func Metrics() gin.HandlerFunc {
	v := setupMetrics()
	return v.Handler()
}

// setupMetrics is the helper function to create the metrics from the CLI context.
func setupMetrics() *metrics.Metrics {
	return metrics.New()
}
