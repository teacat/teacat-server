package main

import (
	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
)

// The functions, structs down below are the core methods,
// you shouldn't edit them until you know what you're doing,
// or you understand how KitSvc works.
//
// Or if you are brave enough ;)

// InstrumentingMiddleware represents a middleware of the instrumenting.
type InstrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	countResult    metrics.Histogram
	Service
}

// createInstruMiddleware creates the instrumenting middleware.
func createInstruMiddleware() ServiceMiddleware {

	fieldKeys := []string{"method", "error"}

	// Number of requests received.
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: viper.GetString("prometheus.namespace"),
		Subsystem: viper.GetString("prometheus.subsystem"),
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)

	// Total duration of requests in microseconds.
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: viper.GetString("prometheus.namespace"),
		Subsystem: viper.GetString("prometheus.subsystem"),
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	// The result of each count method.
	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: viper.GetString("prometheus.namespace"),
		Subsystem: viper.GetString("prometheus.subsystem"),
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{})

	return func(next Service) Service {
		return InstrumentingMiddleware{requestCount, requestLatency, countResult, next}
	}
}
