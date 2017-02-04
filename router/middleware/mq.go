package middleware

import (
	"github.com/TeaMeow/KitSvc/module/mq"
	"github.com/TeaMeow/KitSvc/module/mq/mqstore"
	"github.com/codegangsta/cli"
	"github.com/gin-gonic/gin"
)

func MQ(c *cli.Context) gin.HandlerFunc {
	v := setupMQ(c)
	return func(c *gin.Context) {
		mq.ToContext(c, v)
		c.Next()
	}
}

func setupMQ(c *cli.Context) mq.MQ {
	return mqstore.NewProducer(
		c.String("nsq-producer"),
		c.String("nsq-producer-http"),
		c.StringSlice("nsq-lookupds"),
	)
}
