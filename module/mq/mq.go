package mq

import "github.com/gin-gonic/gin"

var (
	MsgSendMail = "send_mail"
)

// MQ wraps the functions that interactive with the message queue.
type MQ interface {
	Publish(M)
}

type M struct {
	Topic string
	Data  interface{}
}

// UserCreated handles the `user.created` event.
func Publish(c *gin.Context, m M) {
	FromContext(c).Publish(m)
}
