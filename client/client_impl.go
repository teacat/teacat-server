package client

import (
	"fmt"
	"strconv"

	"github.com/TeaMeow/KitSvc/model"
	"github.com/levigross/grequests"
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

// PostUser creates a new user account.
func (c *client) PostUser(in *model.User) (out *model.User, err error) {
	resp, err := grequests.Post(uri(pathUser, c.base), &grequests.RequestOptions{
		JSON: in,
	})
	resp.JSON(&out)
	return
}

// GetUser gets an user by the user identifier.
func (c *client) GetUser(username string) (out *model.User, err error) {
	resp, err := grequests.Get(uri(pathSpecifiedUser, c.base, username), &grequests.RequestOptions{})
	resp.JSON(&out)
	return
}

// PutUser updates an user account information.
func (c *client) PutUser(id int, in *model.User) (out *model.User, err error) {
	resp, err := grequests.Put(uri(pathSpecifiedUser, c.base, id), &grequests.RequestOptions{
		JSON: in,
		Headers: map[string]string{
			"Authorization": "Bearer " + c.token,
		},
	})
	resp.JSON(&out)
	return
}

// DeleteUser deletes the user by the user identifier.
func (c *client) DeleteUser(id int) (err error) {
	_, err = grequests.Delete(uri(pathSpecifiedUser, c.base, id), &grequests.RequestOptions{
		Headers: map[string]string{
			"Authorization": "Bearer " + c.token,
		},
	})
	return
}

// PostToken generates the authentication token
// if the password was matched with the specified account.
func (c *client) PostToken(in *model.User) (out *model.Token, err error) {
	resp, err := grequests.Post(uri(pathAuth, c.base), &grequests.RequestOptions{
		JSON: in,
		Headers: map[string]string{
			"Authorization": "Bearer " + c.token,
		},
	})
	resp.JSON(&out)
	return
}

//
// Helper functions
//

// uri combines the path.
func uri(path string, params ...interface{}) string {
	for i, v := range params {
		switch v.(type) {
		case int:
			params[i] = strconv.Itoa(v.(int))
		}
	}
	return fmt.Sprintf(path, params...)
}
