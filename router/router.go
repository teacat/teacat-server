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

	// Health check handler.
	g.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"pong": "pong",
		})
	})

	// Event handlers.
	e := eventutil.New()
	e.POST("/event-store/user.create/", "user.create", server.CreateUser)

	//e.GET("/kit_metrics", server.)

	return g, e
}
