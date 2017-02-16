package middleware

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/TeaMeow/KitSvc/errno"
	"github.com/TeaMeow/KitSvc/module/logger"
	"github.com/gin-gonic/gin"
	"github.com/willf/pad"
)

// Logging is a middleware function that logs the each request.
func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now().UTC()
		path := c.Request.URL.Path
		// Continue.
		c.Next()
		// Skip for the health check requests.
		if path == "/metrics" || path == "/sd/health" || path == "/sd/ram" || path == "/sd/cpu" || path == "/sd/disk" {
			return
		}
		// Calculates the latency.
		end := time.Now().UTC()
		latency := end.Sub(start)

		// The basic informations.
		status := c.Writer.Status()
		method := c.Request.Method
		ip := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// Create the symbols for each status.
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

		// Data fields that will be recorded into the log files.
		fields := logrus.Fields{
			"user_agent": userAgent,
		}
		// Append the error to the fields so we can record it.
		if len(c.Errors) != 0 {
			for k, v := range c.Errors {
				// Skip if it's the Gin internal error.
				if !v.IsType(gin.ErrorTypePrivate) {
					continue
				}
				// The field name with the `error_INDEX` format.
				errorKey := fmt.Sprintf("error_%d", k)

				switch v.Err.(type) {
				case *errno.Err:
					e := v.Err.(*errno.Err)
					fields[errorKey] = fmt.Sprintf("%s[%s:%d]", e.Code, e.Path, e.Line)
					c.String(e.StatusCode, fmt.Sprintf("%s (Code: %s)", e.Message, e.Code))
				default:
					fields[errorKey] = fmt.Sprintf("%s", v.Err)
				}
			}
		}

		// Example: ● 200 |  102.268592ms |    127.0.0.1 | POST  /user
		logger.InfoFields(fmt.Sprintf("%s | %13s | %12s | %s %s", statusString, latency, ip, pad.Right(method, 5, " "), path), fields)
	}
}
