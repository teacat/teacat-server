package event

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

// Key is the key name of the event in the Gin context.
const Key = "event"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the event store associated with this context.
func FromContext(c context.Context) Event {
	return c.Value(Key).(Event)
}

// ToContext adds the event store to this context if it supports
// the Setter interface.
func ToContext(c Setter, event Event) {
	c.Set(Key, event)
}

// Event wraps the functions that interactive with the Event Store.
type Event interface {
	Send(E)
}

// E represents an event.
type E struct {
	Stream   string
	Data     interface{}
	Metadata map[string]string
}

// Send the event to Event Store.
func Send(c *gin.Context, evt E) {
	FromContext(c).Send(evt)
}
