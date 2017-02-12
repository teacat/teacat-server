package event

import (
	"github.com/TeaMeow/KitSvc/model"
	"github.com/gin-gonic/gin"
)

// Event wraps the functions that interactive with the Event Store.
type Event interface {
	UserCreated(*model.User) error
}

// UserCreated handles the `user.created` event.
func UserCreated(c *gin.Context, user *model.User) error {
	return FromContext(c).UserCreated(user)
}
