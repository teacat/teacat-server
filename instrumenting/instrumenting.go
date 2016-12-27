package instrumenting

import (
	"github.com/TeaMeow/KitSvc/config"
	"github.com/TeaMeow/KitSvc/service"
	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

type Middleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	countResult    metrics.Histogram
	service.Service
}

func CreateMiddleware(c config.Context) service.Middleware {

	fieldKeys := []string{"method", "error"}

	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: c.Prometheus.Namespace,
		Subsystem: c.Prometheus.Subsystem,
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)

	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: c.Prometheus.Namespace,
		Subsystem: c.Prometheus.Subsystem,
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: c.Prometheus.Namespace,
		Subsystem: c.Prometheus.Subsystem,
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{})

	return func(next service.Service) service.Service {
		return Middleware{requestCount, requestLatency, countResult, next}
	}
}
