package mq

import "golang.org/x/net/context"

const Key = "mq"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the Store associated with this context.
func FromContext(c context.Context) MQ {
	return c.Value(Key).(MQ)
}

// ToContext adds the Store to this context if it supports
// the Setter interface.
func ToContext(c Setter, mq MQ) {
	c.Set(Key, mq)
}
