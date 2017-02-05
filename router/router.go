package router

import (
	"os"

	"github.com/TeaMeow/KitSvc/module/metrics"
	"github.com/TeaMeow/KitSvc/module/sd"
	"github.com/TeaMeow/KitSvc/router/middleware/header"
	"github.com/TeaMeow/KitSvc/server"
	"github.com/TeaMeow/KitSvc/shared/eventutil"
	"github.com/TeaMeow/KitSvc/shared/mqutil"
	"github.com/TeaMeow/KitSvc/shared/wsutil"
	"github.com/gin-gonic/gin"
)

//
func Load(g *gin.Engine, e *eventutil.Engine, w *wsutil.Engine, m *mqutil.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// Middlewares.
	g.Use(gin.LoggerWithWriter(os.Stdout, "/metrics", "/sd/health", "/sd/ram", "/sd/cpu", "/sd/disk"))
	g.Use(gin.Recovery())
	g.Use(header.NoCache)
	g.Use(header.Options)
	g.Use(header.Secure)
	g.Use(mw...)

	//TODO: CIRCU BREAKER

	// The common handlers.
	g.POST("/user", server.CreateUser)
	g.GET("/user/:username", server.GetUser)
	g.DELETE("/user/:id", server.DeleteUser)
	g.PUT("/user/:id", server.UpdateUser)
	g.POST("/auth", server.Login)

	// The health check handlers
	// for the service discovery.
	g.GET("/sd/health", sd.HealthCheck)
	g.GET("/sd/disk", sd.DiskCheck)
	g.GET("/sd/cpu", sd.CPUCheck)
	g.GET("/sd/ram", sd.RAMCheck)
	g.GET("/metrics", metrics.PrometheusHandler())

	// Websockets.
	w.Handle("/websocket", server.WebSocket)

	// Message
	m.Capture("user", "send_mail", server.SendMail)

	// The event handlers.
	e.Capture("user.created", server.Created)

	return g
}
