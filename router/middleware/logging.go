package middleware

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/TeaMeow/KitSvc/module/logger"
	"github.com/gin-gonic/gin"
	"github.com/willf/pad"
)

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now().UTC()
		path := c.Request.URL.Path
		c.Next()

		if path == "/metrics" || path == "/sd/health" || path == "/sd/ram" || path == "/sd/cpu" || path == "/sd/disk" {
			return
		}

		end := time.Now().UTC()
		latency := end.Sub(start)

		status := c.Writer.Status()
		method := c.Request.Method
		ip := c.ClientIP()
		userAgent := c.Request.UserAgent()

		logger.InfoFields(fmt.Sprintf("%d | %s | %s | %s %s",
			status,
			pad.Right(latency.String(), 13, " "),
			pad.Right(ip, 12, " "),
			pad.Right(method, 5, " "),
			pad.Right(path, 15, " "),
		), logrus.Fields{
			//"status":     status,
			//"method":     method,
			//"path":       path,
			//"ip":         ip,
			//"duration":   latency,
			"user_agent": userAgent,
		})
	}
}
