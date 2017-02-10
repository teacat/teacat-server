package client

import (
	"github.com/TeaMeow/KitSvc/model"
	"github.com/parnurzeal/gorequest"
)

const (
	pathAuth          = "%s/auth"
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

func (c *client) GetUser(username string) (out *model.User, err []error) {
	_, _, err = gorequest.
		New().
		Get(uri(pathSpecifiedUser, c.base, username)).
		EndStruct(&out)
	return
}

func (c *client) PutUser(id int, in *model.User) (out *model.User, err []error) {
	_, _, err = gorequest.
		New().
		Put(uri(pathSpecifiedUser, c.base, id)).
		Send(in).
		EndStruct(&out)
	return
}

func (c *client) DeleteUser(id int, in *model.User) (out *model.User, err []error) {
	_, _, err = gorequest.
		New().
		Delete(uri(pathSpecifiedUser, c.base, id)).
		Send(in).
		EndStruct(&out)
	return
}

func (c *client) PostAuth(in *model.User) (body string, err []error) {
	_, body, err = gorequest.
		New().
		Post(uri(pathAuth, c.base)).
		Send(in).
		End()
	return
}
