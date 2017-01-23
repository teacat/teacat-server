package event

import (
	"golang.org/x/net/context"
)

const Key = "event"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the Store associated with this context.
func FromContext(c context.Context) Event {
	return c.Value(Key).(Event)
}

// ToContext adds the Store to this context if it supports
// the Setter interface.
func ToContext(c Setter, event Event) {
	c.Set(Key, event)
}
