package middleware

import (
	"github.com/TeaMeow/KitSvc/module/event"
	"github.com/TeaMeow/KitSvc/module/event/eventstore"
	"github.com/TeaMeow/KitSvc/shared/eventutil"
	"github.com/codegangsta/cli"
	"github.com/gin-gonic/gin"
)

func Event(c *cli.Context, e *eventutil.Engine, replayed chan<- bool, deployed <-chan bool) gin.HandlerFunc {
	v := setupEvent(c, e, replayed, deployed)
	return func(c *gin.Context) {
		event.ToContext(c, v)
		c.Next()
	}
}

func setupEvent(c *cli.Context, e *eventutil.Engine, replayed chan<- bool, deployed <-chan bool) event.Event {
	return eventstore.NewClient(
		c.String("url"),
		c.String("es-url"),
		c.String("es-username"),
		c.String("es-password"),
		e,
		replayed,
		deployed,
	)
}
