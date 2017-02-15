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

		statusString := ""
		switch {
		case status >= 500:
			statusString = fmt.Sprintf("▲ %d", status)
		case status >= 400:
			statusString = fmt.Sprintf("▲ %d", status)
		case status >= 300:
			statusString = fmt.Sprintf("■ %d", status)
		case status >= 100:
			statusString = fmt.Sprintf("● %d", status)
		}

		fields := logrus.Fields{
			"user_agent": userAgent,
		}
		if len(c.Errors) != 0 {
			for k, v := range c.Errors {
				if !v.IsType(gin.ErrorTypePrivate) {
					continue
				}

				errorKey := fmt.Sprintf("error_%d", k)

				if v.Meta != nil {
					m := v.Meta.(logger.RouteError)
					fields[errorKey] = fmt.Sprintf("%s[%s:%d]", m.Code, m.Path, m.Line)

				} else {
					fields[errorKey] = fmt.Sprintf("%s", v.Err)
				}

			}
		}
		logger.InfoFields(fmt.Sprintf("%s | %13s | %12s | %s %s", statusString, latency, ip, pad.Right(method, 5, " "), path), fields)

	}
}
