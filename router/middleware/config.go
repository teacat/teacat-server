package middleware

import (
	"github.com/TeaMeow/KitSvc/model"
	"github.com/codegangsta/cli"
	"github.com/gin-gonic/gin"
)

const ConfigKey = "config"

func Config(cli *cli.Context) gin.HandlerFunc {
	v := setupConfig(cli)
	return func(c *gin.Context) {
		c.Set(ConfigKey, v)
	}
}

func setupConfig(c *cli.Context) *model.Config {
	return &model.Config{
		JWTSecret: c.String("jwt-secret"),
	}
}
