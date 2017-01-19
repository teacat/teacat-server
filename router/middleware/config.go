package middleware

import (
	"github.com/TeaMeow/KitSvc/model"
	"github.com/codegangsta/cli"
	"github.com/gin-gonic/gin"
)

const configKey = "config"

func Config(cli *cli.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(configKey, &model.Config{
			JWTSecret: cli.String("jwt-secret"),
		})
	}
}
