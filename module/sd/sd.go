package sd

import (
	"os"
	"os/signal"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/TeaMeow/KitSvc/module/logger"
	"github.com/TeaMeow/KitSvc/version"
	"github.com/codegangsta/cli"
	"github.com/hashicorp/consul/api"
	"github.com/satori/go.uuid"
)

// newClient creates a new Consul api client.
func newClient(c *cli.Context) *api.Client {
	apiConfig := api.DefaultConfig()
	apiClient, err := api.NewClient(apiConfig)
	if err != nil {
		logger.FatalFields("Error occurred while creating the Consul API client.", logrus.Fields{
			"err": err,
		})
	}
	return apiClient
}

// Register the service to the service registry.
func Register(c *cli.Context) {
	client := newClient(c)
	// Create a random id.
	id := uuid.NewV4().String()
	// Append the service version in the consul tags.
	tags := c.StringSlice("consul-tags")
	tags = append(tags, version.Version)

	// The service information.
	info := &api.AgentServiceRegistration{
		ID:   id,
		Name: c.String("name"),
		Port: c.Int("port"),
		Tags: tags,
	}

	if err := client.Agent().ServiceRegister(info); err != nil {
		logger.FatalFields("Error occurred while registering to the service registry (Is consul running?).", logrus.Fields{
			"err": err,
		})
	}
	logger.InfoFields("The service has been registered to the Consul successfully.", logrus.Fields{
		"id":   id,
		"tags": strings.Join(tags, ", "),
	})
	// Register the health check handlers.
	registerChecks(c, client, id)
	// Deregister the service when exiting the program.
	deregister(client, id)
}

// registerChecks register the health check handlers to the service registry.
func registerChecks(c *cli.Context, client *api.Client, id string) {
	checks := []*api.AgentCheckRegistration{
		{
			Name:      "Service Router",
			ServiceID: id,
			AgentServiceCheck: api.AgentServiceCheck{
				HTTP:     c.String("url") + "/sd/health",
				Interval: c.String("consul-check_interval"),
				Timeout:  c.String("consul-check_timeout"),
			},
		},
		{
			Name:      "Disk Usage",
			Notes:     "Critical 5%, warning 10% free",
			ServiceID: id,
			AgentServiceCheck: api.AgentServiceCheck{
				HTTP:     c.String("url") + "/sd/disk",
				Interval: c.String("consul-check_interval"),
				Timeout:  c.String("consul-check_timeout"),
			},
		},
		{
			Name:      "Load Average",
			Notes:     "Critical load average 2, warning load average 1",
			ServiceID: id,
			AgentServiceCheck: api.AgentServiceCheck{
				HTTP:     c.String("url") + "/sd/cpu",
				Interval: c.String("consul-check_interval"),
				Timeout:  c.String("consul-check_timeout"),
			},
		},
		{
			Name:      "RAM Usage",
			Notes:     "Critical 5%, warning 10% free",
			ServiceID: id,
			AgentServiceCheck: api.AgentServiceCheck{
				HTTP:     c.String("url") + "/sd/ram",
				Interval: c.String("consul-check_interval"),
				Timeout:  c.String("consul-check_timeout"),
			},
		},
	}
	for _, v := range checks {
		client.Agent().CheckRegister(v)
	}
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
				logger.FatalFields("Cannot deregister the service from the service registry.", logrus.Fields{
					"err": err,
					"id":  id,
				})
			} else {
				logger.InfoFields("The service has been deregistered from the service registry successfully.", logrus.Fields{
					"id": id,
				})
			}
			os.Exit(1)
		}
	}()
}
