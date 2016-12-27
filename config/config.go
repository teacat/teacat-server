package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Context struct {
	Service
	Database
	Prometheus
	Consul
}

type Service struct {
	Name string
	Port int
}

type Database struct {
	User, Password, Host, Name, Charset, Loc string
	Port                                     int
	ParseTime                                bool
}

type Prometheus struct {
	Namespace, Subsystem string
}

type Consul struct {
	CheckInterval string
	CheckTimeout  string
	Tags          []string
}

func Load() Context {
	viper.SetConfigName("config")
	viper.AddConfigPath("./")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	svc := Service{
		Name: viper.GetString("service.name"),
		Port: viper.GetInt("service.port"),
	}
	db := Database{
		User:      viper.GetString("database.user"),
		Password:  viper.GetString("database.password"),
		Host:      viper.GetString("database.host"),
		Charset:   viper.GetString("database.charset"),
		Loc:       viper.GetString("database.loc"),
		Port:      viper.GetInt("database.port"),
		ParseTime: viper.GetBool("database.parse_time"),
	}
	pt := Prometheus{
		Namespace: viper.GetString("prometheus.namespace"),
		Subsystem: viper.GetString("prometheus.subsystem"),
	}
	cl := Consul{
		CheckInterval: viper.GetString("consul.check.interval"),
		CheckTimeout:  viper.GetString("consul.check.timeout"),
		Tags:          viper.GetStringSlice("consul.tags"),
	}

	return Context{Service: svc, Database: db, Prometheus: pt, Consul: cl}
}
