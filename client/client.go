package client

import (
	"fmt"

	"github.com/TeaMeow/KitSvc/model"
)

type Client interface {
	PostUser(*model.User) (*model.User, []error)
	GetUser(string) (*model.User, []error)
	PutUser(int, *model.User) (*model.User, []error)
	DeleteUser(int, *model.User) (*model.User, []error)
	PostAuth(*model.User) (string, []error)
}

func uri(path string, params ...interface{}) string {
	return fmt.Sprintf(path, params...)
}
