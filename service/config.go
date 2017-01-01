package main

import (
	"fmt"

	"github.com/spf13/viper"
)

func loadConfig(path string) {
	viper.SetConfigName("config")
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
