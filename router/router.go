package router

import (
	"github.com/TeaMeow/KitSvc/module/metrics"
	"github.com/TeaMeow/KitSvc/module/sd"
	"github.com/TeaMeow/KitSvc/router/middleware/header"
	"github.com/TeaMeow/KitSvc/server"
	"github.com/TeaMeow/KitSvc/shared/eventutil"
	"github.com/gin-gonic/gin"
)

func Load(g *gin.Engine, e *eventutil.Engine, m *metrics.Metrics, mw ...gin.HandlerFunc) *gin.Engine {
	// Middlewares.
	g.Use(gin.Logger())
	g.Use(gin.Recovery())
	g.Use(header.NoCache)
	g.Use(header.Options)
	g.Use(header.Secure)
	g.Use(mw...)
	// Prometheus middleware.
	//p := ginprometheus.NewPrometheus("gin")
	//p.Use(g)

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
	g.GET("/metrics", m.Handler())

	// The event handlers.
	e.POST("/es/user.create/", "user.create", server.Created)

	return g
}
