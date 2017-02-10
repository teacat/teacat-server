package main

import (
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/TeaMeow/KitSvc/client"
	"github.com/TeaMeow/KitSvc/model"
	"github.com/TeaMeow/KitSvc/version"
	"github.com/codegangsta/cli"
	"github.com/stretchr/testify/assert"
)

/*PostUser(*model.User) (*model.User, []error)
GetUser(string) (*model.User, []error)
PutUser(int, *model.User) (*model.User, []error)
DeleteUser(int, *model.User) (*model.User, []error)
PostAuth(*model.User) (string, []error)*/

//
var started = make(chan bool)
var c = client.NewClient("http://127.0.0.1:8080")

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

func printErrors(e []error) {
	if len(e) != 0 {
		for _, v := range e {
			logrus.Error(v.Error())
		}
	}
}

func TestPostUser(t *testing.T) {
	assert := assert.New(t)

	u, errs := c.PostUser(&model.User{
		Username: "admin",
		Password: "testtest",
	})
	printErrors(errs)

	assert.Equal(&model.User{
		ID:       1,
		Username: "admin",
		Password: "testtest",
	}, u, "They should be equal.")
}

func TestGetUser(t *testing.T) {
	//assert := assert.New(t)
}

func TestPutUser(t *testing.T) {
	//assert := assert.New(t)
}

func TestDeleteUser(t *testing.T) {
	//assert := assert.New(t)
}

func TestPostAuth(t *testing.T) {
	//assert := assert.New(t)
}
