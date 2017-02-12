package mq

import "golang.org/x/net/context"

// Key is the key name of the message queue in the Gin context.
const Key = "mq"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the message queue associated with this context.
func FromContext(c context.Context) MQ {
	return c.Value(Key).(MQ)
}

// ToContext adds the message queue to this context if it supports
// the Setter interface.
func ToContext(c Setter, mq MQ) {
	c.Set(Key, mq)
}
