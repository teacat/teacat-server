package sd

import (
	"os"
	"os/signal"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/hashicorp/consul/api"
	"github.com/satori/go.uuid"
)

func Wait(c *cli.Context, played <-chan bool) {
	// Block until the events were all replayed.
	<-played

	logrus.Infoln("The events were all replayed, trying to register to the server registry.")

	client := newClient(c)
	register(c, client)

	logrus.Infoln("The service has been registered to the server registry successfully.")
}

func newClient(c *cli.Context) *api.Client {
	apiConfig := api.DefaultConfig()
	apiClient, err := api.NewClient(apiConfig)
	if err != nil {
		logrus.Errorln(err)
		logrus.Fatalln("Error occurred while creating the Consul api client.")
	}
	return apiClient
}

// register register the service to the service registry.
func register(c *cli.Context, client *api.Client) {
	id := uuid.NewV4().String()

	// The information of the health check.

	// The service information.
	info := &api.AgentServiceRegistration{
		ID:   id,
		Name: c.String("name"),
		Port: c.Int("port"),
		Tags: c.StringSlice("consul-tags"),
	}

	// Register the service to the service registry.
	if err := client.Agent().ServiceRegister(info); err != nil {
		logrus.Errorln(err)
		logrus.Fatalln("Error occurred while registering to the service registry (Is consul running?).")
	}

	check := &api.AgentCheckRegistration{
		Name:      "Web",
		Notes:     "Wow",
		ServiceID: id,
		AgentServiceCheck: api.AgentServiceCheck{
			HTTP:     c.String("url") + "/health",
			Interval: c.String("consul-check_interval"),
			Timeout:  c.String("consul-check_timeout"),
		},
	}

	//

	client.Agent().CheckRegister(check)

	//

	check2 := &api.AgentCheckRegistration{
		Name:      "Disk",
		Notes:     "Wow",
		ServiceID: id,
		AgentServiceCheck: api.AgentServiceCheck{
			HTTP:     c.String("url") + "/disk",
			Notes:    "Critical 5%, warning 10% free",
			Interval: c.String("consul-check_interval"),
			Timeout:  c.String("consul-check_timeout"),
		},
	}

	client.Agent().CheckRegister(check2)

	// Deregister the service when exiting the program.
	deregister(client, id)
}

// deregister watching the system signal, deregister the service from the service registry
// when the exit signal was captured.
func deregister(client *api.Client, id string) {
	// Capture the program exit signal.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	go func() {
		for range ch {
			if err := client.Agent().ServiceDeregister(id); err != nil {
				logrus.Errorln(err)
				logrus.Fatalln("Cannot deregister the service from the service registry.")
			} else {
				logrus.Infoln("The service has been deregistered from the service registry successfully.")
			}
			os.Exit(1)
		}
	}()
}
