package main

import (
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/go-kit/kit/log"
	consulsd "github.com/go-kit/kit/sd/consul"
	consulapi "github.com/hashicorp/consul/api"
)

// registerService register the service to the service discovery server(consul).
func registerService(logger log.Logger) {
	p, _ := strconv.Atoi(os.Getenv("KITSVC_PORT"))

	info := consulapi.AgentServiceRegistration{
		Name: os.Getenv("KITSVC_NAME"),
		Port: p,
		Tags: strings.Split(os.Getenv("KITSVC_CONSUL_TAGS"), ","),
		Check: &consulapi.AgentServiceCheck{
			HTTP:     os.Getenv("KITSVC_URL") + "/sd_health",
			Interval: os.Getenv("KITSVC_CONSUL_CHECK_INTERVAL"),
			Timeout:  os.Getenv("KITSVC_CONSUL_CHECK_TIMEOUT"),
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
