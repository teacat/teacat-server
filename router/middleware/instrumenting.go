package middleware

import (
	"github.com/TeaMeow/KitSvc/module/metrics"
	"github.com/gin-gonic/gin"
)

func Metrics() gin.HandlerFunc {
	v := setupMetrics()

	return v.Handler()
}

func setupMetrics() *metrics.Metrics {
	return metrics.New()
}
