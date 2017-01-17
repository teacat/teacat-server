package middleware

import (
	"github.com/codegangsta/cli"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

func Instrumenting(c *cli.Context) gin.HandlerFunc {

	fieldKeys := []string{"method", "uri", "status_code"}

	count := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: c.String("prometheus-namespace"),
			Subsystem: c.String("prometheus-subsystem"),
			Name:      "request_count",
			Help:      "Number of requests received.",
		},
		fieldKeys,
	)
	latency := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: c.String("prometheus-namespace"),
			Subsystem: c.String("prometheus-subsystem"),
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		},
		fieldKeys,
	)

	prometheus.MustRegister(count)
	prometheus.MustRegister(latency)

	return func(c *gin.Context) {

		/*defer func(begin time.Time) {
			method := c.Request.Method
			uri := c.Request.RequestURI
			code := c.Request.Response.StatusCode

			//prometheus.Labels{"method": method, "uri": uri, "status_code": code}
			//logrus.Println(time.Since(begin).Seconds())
		}(time.Now())
		c.Next()*/
	}
}
