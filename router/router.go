package router

import (
	"net/http"

	"github.com/TeaMeow/KitSvc/server"
	"github.com/TeaMeow/KitSvc/shared/eventutil"
	"github.com/gin-gonic/gin"
)

func Load(middleware ...gin.HandlerFunc) (http.Handler, *eventutil.Engine) {

	// Gin engine and middlewares.
	g := gin.Default()
	g.Use(gin.Recovery())
	g.Use(middleware...)

	// Command routes.
	g.POST("/user", server.CreateUser)
	g.GET("/user/:id", server.GetUser)
	g.DELETE("/user/:id", server.DeleteUser)
	g.PUT("/user/:id", server.UpdateUser)

	// Event handlers.
	e := eventutil.New()
	e.POST("/event_store/user.create/", "user.create", server.CreateUser)

	// Service handlers.
	//e.GET("/kit_sd_health", server.)
	//e.GET("/kit_metrics", server.)

	return g, e
}
