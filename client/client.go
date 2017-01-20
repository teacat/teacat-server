package client

import (
	"github.com/TeaMeow/KitSvc/model"
	"github.com/TeaMeow/KitSvc/protobuf"
)

type Client interface {
	PostUser(*model.User) protobuf.CreateUserResponse
	GetUser()
	PutUser()
	DeleteUser()
	PostAuth()
}
