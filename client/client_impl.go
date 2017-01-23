package client

import (
	"github.com/TeaMeow/KitSvc/model"
	"github.com/parnurzeal/gorequest"
)

const (
	pathUser          = "%s/user"
	pathSpecifiedUser = "%s/user/%s"
)

type client struct {
	token string
	base  string
}

func NewClient(uri string) Client {
	return &client{base: uri}
}

func (c *client) PostUser(in *model.User) (out *model.User, err []error) {
	_, _, err = gorequest.
		New().
		Post(uri(pathUser, c.base)).
		Send(in).
		EndStruct(&out)
	return
}

func (c *client) GetUser() {

}

func (c *client) PutUser() {

}

func (c *client) DeleteUser() {

}

func (c *client) PostAuth() {

}
