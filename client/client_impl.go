package client

import (
	"encoding/json"
	"errors"

	"github.com/TeaMeow/KitSvc/model"
	"github.com/parnurzeal/gorequest"
)

const (
	pathAuth          = "%s/user/token"
	pathUser          = "%s/user"
	pathSpecifiedUser = "%s/user/%s"
)

// client represents the http client.
type client struct {
	token string
	base  string
}

// NewClient returns a client at the specified url.
func NewClient(uri string) Client {
	return &client{base: uri}
}

// NewClientToken returns a client at the specified url that authenticates all
// outbound requests with the given token.
func NewClientToken(uri, token string) Client {
	return &client{token: token, base: uri}
}

// handler handles the request error.
func handler(resp gorequest.Response, b string, e []error, out interface{}) (errs []error) {
	// Append the current errors to the error array.
	errs = append(errs, e...)
	// Unmarshal the JSON body to the struct if the response status code is lower than 300 (== OK).
	if resp != nil && resp.StatusCode < 300 {
		// If error occurred while processing the JSON body, append it to the error array.
		if err := json.Unmarshal([]byte(b), &out); err != nil {
			errs = append(errs, err)
		}
		return
	}
	// Append the response body to the error array
	// if the response is not in the "OK range" (we assumed the response is the error message.).
	errs = append(errs, errors.New(b))
	return
}

// PostUser creates a new user account.
func (c *client) PostUser(in *model.User) (out *model.User, errs []error) {
	resp, b, e := gorequest.
		New().
		Post(uri(pathUser, c.base)).
		Send(in).
		End()

	errs = handler(resp, b, e, &out)
	return
}

// GetUser gets an user by the user identifier.
func (c *client) GetUser(username string) (out *model.User, errs []error) {
	resp, b, e := gorequest.
		New().
		Get(uri(pathSpecifiedUser, c.base, username)).
		End()

	errs = handler(resp, b, e, &out)
	return
}

// PutUser updates an user account information.
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

// DeleteUser deletes the user by the user identifier.
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

// PostToken generates the authentication token
// if the password was matched with the specified account.
func (c *client) PostToken(in *model.User) (out *model.Token, errs []error) {
	resp, b, e := gorequest.
		New().
		Post(uri(pathAuth, c.base)).
		Send(in).
		End()
	errs = handler(resp, b, e, &out)
	return
}
