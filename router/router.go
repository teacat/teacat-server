package router

import (
	"net/http"

	"github.com/TeaMeow/KitSvc/module/event"
	"github.com/TeaMeow/KitSvc/module/metrics"
	"github.com/TeaMeow/KitSvc/module/mq"
	"github.com/TeaMeow/KitSvc/module/sd"
	"github.com/TeaMeow/KitSvc/router/middleware/header"
	"github.com/TeaMeow/KitSvc/service"
	"github.com/TeaMeow/KitSvc/shared/eventutil"
	"github.com/TeaMeow/KitSvc/shared/mqutil"
	"github.com/TeaMeow/KitSvc/shared/wsutil"
	"github.com/gin-gonic/gin"
)

// Load loads the middlewares, routes, handlers.
func Load(g *gin.Engine, e *eventutil.Engine, w *wsutil.Engine, m *mqutil.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// Middlewares.
	g.Use(gin.Recovery())
	g.Use(header.NoCache)
	g.Use(header.Options)
	g.Use(header.Secure)
	g.Use(mw...)
	// 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route.")
	})

	// The common handlers.
	user := g.Group("/user")
	{
		user.POST("", service.CreateUser)
		user.GET("/:username", service.GetUser)
		user.DELETE("/:id", service.DeleteUser)
		user.PUT("/:id", service.UpdateUser)
		user.POST("/token", service.PostToken)
	}

	// The health check handlers
	// for the service discovery.
	svcd := g.Group("/sd")
	{
		svcd.GET("/health", sd.HealthCheck)
		svcd.GET("/disk", sd.DiskCheck)
		svcd.GET("/cpu", sd.CPUCheck)
		svcd.GET("/ram", sd.RAMCheck)
	}

	// Prometheus metrics handler.
	g.GET("/metrics", metrics.PrometheusHandler())

	// WebSockets.
	w.Handle("/websocket", service.WatchUser)

	// Message handlers.
	m.Capture("user", mq.MsgSendMail, service.SendMail)

	// Event handlers.
	e.Capture(event.EvtUserCreated, service.UserCreated)

	return g
}
