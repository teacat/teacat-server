package middleware

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/TeaMeow/KitSvc/module/logger"
	"github.com/fatih/color"
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

		statusString := ""
		switch {
		case status >= 500:
			statusString = color.New(color.BgRed).Add(color.FgWhite).Sprintf(" %d ", status)
		case status >= 400:
			statusString = color.New(color.BgRed).Add(color.FgWhite).Sprintf(" %d ", status)
		case status >= 300:
			statusString = color.New(color.BgYellow).Add(color.FgBlack).Sprintf(" %d ", status)
		case status >= 100:
			statusString = color.New(color.BgGreen).Add(color.FgWhite).Sprintf(" %d ", status)
		}

		logger.InfoFields(fmt.Sprintf("%s | %s | %s | %s %s",
			statusString,
			pad.Right(latency.String(), 13, " "),
			pad.Right(ip, 12, " "),
			pad.Right(method, 5, " "),
			path,
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
