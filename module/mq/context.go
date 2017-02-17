package mq

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

// Key is the key name of the message queue in the Gin context.
const Key = "mq"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// MQ wraps the functions that interactive with the message queue.
type MQ interface {
	Publish(M)
}

// M represents a message.
type M struct {
	Topic string
	Data  interface{}
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

// Publish the message.
func Publish(c *gin.Context, m M) {
	FromContext(c).Publish(m)
}
