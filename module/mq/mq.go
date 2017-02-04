package mq

import (
	"github.com/TeaMeow/KitSvc/model"
	"github.com/gin-gonic/gin"
)

type MQ interface {
	SendMail(*model.User) error
}

func SendMail(c *gin.Context, user *model.User) error {
	return FromContext(c).SendMail(user)
}
