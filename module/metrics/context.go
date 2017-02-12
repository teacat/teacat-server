package metrics

import (
	"golang.org/x/net/context"
)

// Key is the key name of the metrics in the Gin context.
const Key = "metrics"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the metric associated with this context.
func FromContext(c context.Context) *Metrics {
	return c.Value(Key).(*Metrics)
}

// ToContext adds the metric to this context if it supports
// the Setter interface.
func ToContext(c Setter, metrics *Metrics) {
	c.Set(Key, metrics)
}
