package event

import "github.com/gin-gonic/gin"

var (
	EvtUserCreated = "user_created"
	EvtUserDeleted = "user_deleted"
)

// Event wraps the functions that interactive with the Event Store.
type Event interface {
	Send(E)
}

type E struct {
	Stream   string
	Data     interface{}
	Metadata map[string]string
}

// UserCreated handles the `user.created` event.
func Send(c *gin.Context, evt E) {
	FromContext(c).Send(evt)
}
