package main

import (
	"os"
	"testing"

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
var serverReady = make(chan bool)

func TestMain(t *testing.T) {
	app := cli.NewApp()
	app.Name = "service"
	app.Version = version.Version
	app.Usage = "starts the service daemon."
	app.Action = func(c *cli.Context) {
		server(c, serverReady)
	}
	app.Flags = serverFlags

	go app.Run(os.Args)
	<-serverReady
}

func TestPostUser(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(123, 123, "they should be equal")
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
