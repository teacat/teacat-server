package discovery

import (
	"strconv"

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
			HTTP:     "http://localhost:" + strconv.Itoa(c.Service.Port) + "/health",
			Interval: c.Consul.CheckInterval,
			Timeout:  c.Consul.CheckTimeout,
		},
	}
	// DEREGISTRE
	// DDDDDD
	// DDD
	apiConfig := consulapi.DefaultConfig()
	apiClient, _ := consulapi.NewClient(apiConfig)
	client := consulsd.NewClient(apiClient)
	reg := consulsd.NewRegistrar(client, &info, logger)
	reg.Register()
}
