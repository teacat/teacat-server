package middleware

import (
	"github.com/TeaMeow/KitSvc/module/mq"
	"github.com/TeaMeow/KitSvc/module/mq/mqstore"
	"github.com/TeaMeow/KitSvc/shared/mqutil"
	"github.com/codegangsta/cli"
	"github.com/gin-gonic/gin"
)

// MQ is a middleware function that initializes the message queue and attaches to
// the context of every request context.
func MQ(c *cli.Context, m *mqutil.Engine, deployed <-chan bool) gin.HandlerFunc {
	v := setupMQ(c, m, deployed)
	return func(c *gin.Context) {
		mq.ToContext(c, v)
		c.Next()
	}
}

// setupMQ is the helper function to create the message queue from the CLI context.
func setupMQ(c *cli.Context, m *mqutil.Engine, deployed <-chan bool) mq.MQ {
	return mqstore.NewProducer(
		c.String("url"),
		c.String("nsq-producer"),
		c.String("nsq-producer-http"),
		c.StringSlice("nsq-lookupds"),
		m,
		deployed,
	)
}
