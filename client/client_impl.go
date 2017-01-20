package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/TeaMeow/KitSvc/model"
	"github.com/TeaMeow/KitSvc/protobuf"
	"github.com/gogo/protobuf/proto"
)

const (
	pathUser          = "%s/user"
	pathSpecifiedUser = "%s/user/%s"
)

type client struct {
	client *http.Client
	token  string
	base   string
}

func NewClient(uri string) Client {
	return &client{client: http.DefaultClient, base: uri}
}

func (c *client) PostUser(in *model.User) protobuf.CreateUserResponse {
	pb := protobuf.CreateUserRequest{
		Username: in.Username,
		Password: in.Password,
	}
	mIn, err := proto.Marshal(&pb)
	if err != nil {
		panic(err)
	}
	uri := fmt.Sprintf(pathUser, c.base)

	// Request.
	req, err := http.NewRequest("POST", uri, bytes.NewReader(mIn))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/x-protobuf")

	res, err := c.client.Do(req)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	// Response to proto buffer.
	var out protobuf.CreateUserResponse
	proto.Unmarshal(body, &out)

	return out
}

func (c *client) GetUser() {

}

func (c *client) PutUser() {

}

func (c *client) DeleteUser() {

}

func (c *client) PostAuth() {

}
