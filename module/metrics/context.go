package metrics

import (
	"golang.org/x/net/context"
)

const Key = "metrics"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the Store associated with this context.
func FromContext(c context.Context) *Metrics {
	return c.Value(Key).(*Metrics)
}

// ToContext adds the Store to this context if it supports
// the Setter interface.
func ToContext(c Setter, metrics *Metrics) {
	c.Set(Key, metrics)
}
