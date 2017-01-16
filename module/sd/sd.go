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

	client := newClient(c)
	register(c, client)
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
	check := &api.AgentServiceCheck{
		HTTP:     c.String("url") + "/health",
		Interval: c.String("consul-check_interval"),
		Timeout:  c.String("consul-check_timeout"),
	}

	// The service information.
	info := &api.AgentServiceRegistration{
		ID:    id,
		Name:  c.String("name"),
		Port:  c.Int("port"),
		Tags:  c.StringSlice("consul-tags"),
		Check: check,
	}

	// Register the service to the service registry.
	client.Agent().ServiceRegister(info)

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
			client.Agent().ServiceDeregister(id)
			os.Exit(1)
		}
	}()
}
