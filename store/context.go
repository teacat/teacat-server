package store

import (
	"golang.org/x/net/context"
)

// Key is the key name of the store in the Gin context.
const Key = "store"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the Store associated with this context.
func FromContext(c context.Context) Store {
	return c.Value(Key).(Store)
}

// ToContext adds the Store to this context if it supports
// the Setter interface.
func ToContext(c Setter, store Store) {
	c.Set(Key, store)
}
