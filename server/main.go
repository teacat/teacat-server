package main

import (
	"os"

	"../version"
	"github.com/codegangsta/cli"
)

func main() {
	started := make(chan bool)

	app := cli.NewApp()
	app.Name = "service"
	app.Version = version.Version
	app.Usage = "starts the service daemon."
	app.Action = func(c *cli.Context) {
		server(c, started)
	}
	app.Flags = serverFlags

	app.Run(os.Args)
}
