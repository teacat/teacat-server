package middleware

import (
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := logrus.StandardLogger()
		start := time.Now().UTC()
		path := c.Request.URL.Path
		c.Next()
		end := time.Now().UTC()
		latency := end.Sub(start)

		formatter := &logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02-15:04:05",
		}
		logrus.SetFormatter(formatter)
		//logrus.SetFormatter(&logrus.JSONFormatter{})
		file, err := os.OpenFile("./logrus.log", os.O_APPEND|os.O_WRONLY, 0666)
		if err == nil {
			logger.Out = file
		} else {
			logger.Info("Failed to log to file, using default stderr")
		}

		logger.WithFields(logrus.Fields{
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       path,
			"ip":         c.ClientIP(),
			"duration":   latency,
			"user_agent": c.Request.UserAgent(),
		}).Info()
	}
}
