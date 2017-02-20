package main

import (
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/TeaMeow/KitSvc/client"
	"github.com/TeaMeow/KitSvc/model"
	"github.com/TeaMeow/KitSvc/shared/token"
	"github.com/TeaMeow/KitSvc/version"
	"github.com/codegangsta/cli"
	"github.com/stretchr/testify/assert"
)

//
var started = make(chan bool)
var c = client.NewClient("http://127.0.0.1:8080")
var ct client.Client

func printErrors(e []error) {
	if len(e) != 0 {
		for _, v := range e {
			logrus.Error(v.Error())
		}
	}
}

func TestMain(t *testing.T) {
	app := cli.NewApp()
	app.Name = "service"
	app.Version = version.Version
	app.Usage = "starts the service daemon."
	app.Action = func(c *cli.Context) {
		server(c, started)
	}
	app.Flags = serverFlags

	go app.Run(os.Args)
	<-started
}

func TestPostUser(t *testing.T) {
	assert := assert.New(t)

	u, err := c.PostUser(&model.User{
		Username: "admin",
		Password: "testtest",
	})
	assert.True(err == nil)
	err = u.Compare("testtest")
	assert.True(err == nil)
}

func TestGetUser(t *testing.T) {
	assert := assert.New(t)

	u, err := c.GetUser("admin")
	assert.True(err == nil)

	err = u.Compare("testtest")
	assert.True(err == nil)
}

func TestPostToken(t *testing.T) {
	assert := assert.New(t)

	tkn, err := c.PostToken(&model.User{
		Username: "admin",
		Password: "testtest",
	})
	assert.True(err == nil)

	ctx, err := token.Parse(tkn.Token, "4Rtg8BPKwixXy2ktDPxoMMAhRzmo9mmuZjvKONGPZZQSaJWNLijxR42qRgq0iBb5")
	assert.True(err == nil)
	assert.Equal(&token.Context{
		ID:       1,
		Username: "admin",
	}, ctx, "They should be equal.")

	ct = client.NewClientToken("http://127.0.0.1:8080", tkn.Token)
}

func TestPutUser(t *testing.T) {
	assert := assert.New(t)

	u, err := ct.PutUser(1, &model.User{
		Username: "admin",
		Password: "newpassword",
	})
	assert.True(err == nil)

	err = u.Compare("newpassword")
	assert.True(err == nil, "They should be match.")
}

func TestDeleteUser(t *testing.T) {
	assert := assert.New(t)

	err := ct.DeleteUser(1)
	assert.True(err == nil)

	u, err := c.GetUser("admin")
	assert.True(err == nil)
	assert.Empty(u)
}
