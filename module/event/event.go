package event

import (
	"github.com/TeaMeow/KitSvc/model"
	"github.com/gin-gonic/gin"
)

type Event interface {
	UserCreated(*model.User) error
}

func UserCreated(c *gin.Context, user *model.User) error {
	return FromContext(c).UserCreated(user)
}
