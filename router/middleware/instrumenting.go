package middleware

import (
	"github.com/TeaMeow/KitSvc/module/metrics"
	"github.com/gin-gonic/gin"
)

func Metrics() gin.HandlerFunc {
	v := setupMetrics()

	return v.Handler()
	//return func(c *gin.Context) {
	//	metrics.ToContext(c, v)
	//    return metrics(v)
	//	//c.Next()
	//}
}

func setupMetrics() *metrics.Metrics {
	return metrics.New()
}
