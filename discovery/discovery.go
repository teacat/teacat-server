package discovery

import (
	"os"
	"os/signal"

	"github.com/TeaMeow/KitSvc/config"
	"github.com/go-kit/kit/log"
	consulsd "github.com/go-kit/kit/sd/consul"
	consulapi "github.com/hashicorp/consul/api"
)

func Register(c *config.Context, logger log.Logger) {

	info := consulapi.AgentServiceRegistration{
		Name: c.Service.Name,
		Port: c.Service.Port,
		Tags: c.Consul.Tags,
		Check: &consulapi.AgentServiceCheck{
			HTTP:     c.Service.URL + "/health",
			Interval: c.Consul.CheckInterval,
			Timeout:  c.Consul.CheckTimeout,
		},
	}

	apiConfig := consulapi.DefaultConfig()
	apiClient, _ := consulapi.NewClient(apiConfig)
	client := consulsd.NewClient(apiClient)
	reg := consulsd.NewRegistrar(client, &info, logger)

	// Deregister the service when ctrl+c
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		for range ch {
			reg.Deregister()
			os.Exit(1)
		}
	}()

	// Register the service
	reg.Register()
}
