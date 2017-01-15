package router

import (
	"net/http"

	"github.com/TeaMeow/KitSvc/server"
	"github.com/gin-gonic/gin"
)

func Load(middleware ...gin.HandlerFunc) http.Handler {
	e := gin.New()
	e.Use(middleware...)

	e.POST("/user", server.CreateUser)
	e.GET("/user/:id", server.GetUser)
	e.DELETE("/user/:id", server.DeleteUser)
	e.PUT("/user/:id", server.UpdateUser)

	return e
}
