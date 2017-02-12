package client

import (
	"fmt"
	"strconv"

	"github.com/TeaMeow/KitSvc/model"
)

type Client interface {
	PostUser(*model.User) (*model.User, []error)
	GetUser(string) (*model.User, []error)
	PutUser(int, *model.User) (*model.User, []error)
	DeleteUser(int, *model.User) (*model.User, []error)
	PostToken(*model.User) (*model.Token, []error)
}

func uri(path string, params ...interface{}) string {
	for i, v := range params {
		switch v.(type) {
		case int:
			params[i] = strconv.Itoa(v.(int))
		}
	}

	return fmt.Sprintf(path, params...)
}
