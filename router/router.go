package router

import (
	"net/http"

	"github.com/TeaMeow/KitSvc/module/sd"
	"github.com/TeaMeow/KitSvc/router/middleware"
	"github.com/TeaMeow/KitSvc/server"
	"github.com/TeaMeow/KitSvc/shared/eventutil"
	"github.com/codegangsta/cli"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/mcuadros/go-gin-prometheus"
)

func Load(c *cli.Context, mw ...gin.HandlerFunc) (http.Handler, *eventutil.Engine) {

	// The middlewares.
	g := gin.New()

	g.Use(gin.Logger())
	g.Use(gin.Recovery())
	g.Use(mw...)

	p := ginprometheus.NewPrometheus("gin")
	p.Use(g)

	// The common handlers.
	g.POST("/user", middleware.JWT(c), server.CreateUser)
	g.GET("/user/:id", server.GetUser)
	g.DELETE("/user/:id", server.DeleteUser)
	g.PUT("/user/:id", server.UpdateUser)

	// The health check handlers
	// for the service discovery.
	g.GET("/sd/health", sd.HealthCheck)
	g.GET("/sd/disk", sd.DiskCheck)
	g.GET("/sd/cpu", sd.CPUCheck)
	g.GET("/sd/ram", sd.RAMCheck)

	// The event handlers.
	e := eventutil.New(g)
	e.POST("/es/user.create/", "user.create", server.CreateUser)

	return e.Gin, e
}
