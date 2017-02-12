package middleware

import (
	"github.com/codegangsta/cli"
	"github.com/gin-gonic/gin"
)

// configKey is the key name of the config context in the Gin context.
const configKey = "config"

// Config is a middleware function that initializes the config and attaches to
// the context of every request context.
func Config(cli *cli.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(configKey, cli)
	}
}

// ConfigContext returns the CLI context associated with this context.
func ConfigContext(c *gin.Context) (ctx *cli.Context) {
	conf, _ := c.Get(configKey)
	ctx = conf.(*cli.Context)
	return
}
