package client

import (
	"encoding/json"
	"errors"

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

func NewClientToken(uri, token string) Client {
	return &client{token: token, base: uri}
}

func handler(resp gorequest.Response, b string, e []error, out interface{}) (errs []error) {
	errs = append(errs, e...)
	if resp != nil && resp.StatusCode < 300 {
		if err := json.Unmarshal([]byte(b), &out); err != nil {
			errs = append(errs, err)
		}
		return
	}
	errs = append(errs, errors.New(b))

	return
}

func (c *client) PostUser(in *model.User) (out *model.User, errs []error) {
	resp, b, e := gorequest.
		New().
		Post(uri(pathUser, c.base)).
		Send(in).
		End()

	errs = handler(resp, b, e, &out)
	return
}

func (c *client) GetUser(username string) (out *model.User, errs []error) {
	resp, b, e := gorequest.
		New().
		Get(uri(pathSpecifiedUser, c.base, username)).
		End()

	errs = handler(resp, b, e, &out)
	return
}

func (c *client) PutUser(id int, in *model.User) (out *model.User, errs []error) {
	resp, b, e := gorequest.
		New().
		Put(uri(pathSpecifiedUser, c.base, id)).
		Set("Authorization", "Bearer "+c.token).
		Send(in).
		End()

	errs = handler(resp, b, e, &out)
	return
}

func (c *client) DeleteUser(id int, in *model.User) (out *model.User, errs []error) {
	resp, b, e := gorequest.
		New().
		Set("Authorization", "Bearer "+c.token).
		Delete(uri(pathSpecifiedUser, c.base, id)).
		Send(in).
		End()

	errs = handler(resp, b, e, &out)
	return
}

func (c *client) PostAuth(in *model.User) (out *model.Token, errs []error) {
	resp, b, e := gorequest.
		New().
		Post(uri(pathAuth, c.base)).
		Send(in).
		End()
	errs = handler(resp, b, e, &out)
	return
}
