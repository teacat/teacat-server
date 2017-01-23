package client

import (
	"fmt"

	"github.com/TeaMeow/KitSvc/model"
)

type Client interface {
	PostUser(*model.User) (out *model.User, err []error)
	GetUser()
	PutUser()
	DeleteUser()
	PostAuth()
}

func uri(path string, base string) string {
	return fmt.Sprintf(path, base)
}
