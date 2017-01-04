package main

import (
	"os"
	"os/signal"

	"github.com/go-kit/kit/log"
	consulsd "github.com/go-kit/kit/sd/consul"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
)

// registerService register the service to the service discovery server(consul).
func registerService(logger log.Logger) {

	info := consulapi.AgentServiceRegistration{
		Name: viper.GetString("service.name"),
		Port: viper.GetInt("service.port"),
		Tags: viper.GetStringSlice("service.tags"),
		Check: &consulapi.AgentServiceCheck{
			HTTP:     viper.GetString("service.url") + "/health",
			Interval: viper.GetString("consul.check.interval"),
			Timeout:  viper.GetString("consul.check.timeout"),
		},
	}

	apiConfig := consulapi.DefaultConfig()
	apiClient, _ := consulapi.NewClient(apiConfig)
	client := consulsd.NewClient(apiClient)
	reg := consulsd.NewRegistrar(client, &info, logger)

	// Deregister the service when exiting the program.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		for range ch {
			reg.Deregister()
			os.Exit(1)
		}
	}()

	// Register the service.
	reg.Register()
}
