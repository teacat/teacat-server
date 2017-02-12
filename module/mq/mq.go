package mq

import (
	"github.com/TeaMeow/KitSvc/model"
	"github.com/gin-gonic/gin"
)

// MQ wraps the functions that interactive with the message queue.
type MQ interface {
	SendMail(*model.User) error
}

func SendMail(c *gin.Context, user *model.User) error {
	return FromContext(c).SendMail(user)
}
