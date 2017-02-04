package router

import (
	"net/http"
	"os"
	"time"

	"github.com/TeaMeow/KitSvc/module/metrics"
	"github.com/TeaMeow/KitSvc/module/sd"
	"github.com/TeaMeow/KitSvc/router/middleware/header"
	"github.com/TeaMeow/KitSvc/server"
	"github.com/TeaMeow/KitSvc/shared/eventutil"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

type route struct {
	http      string
	websocket string
	topic     string
	channel   string
	event     string
	handler   gin.HandlerFunc
}

//
func Load(g *gin.Engine, e *eventutil.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// Middlewares.
	g.Use(gin.LoggerWithWriter(os.Stdout, "/metrics", "/sd/health", "/sd/ram", "/sd/cpu", "/sd/disk"))
	g.Use(gin.Recovery())
	g.Use(header.NoCache)
	g.Use(header.Options)
	g.Use(header.Secure)
	g.Use(mw...)

	routes := []route{
		{
			http:    "POST /user",
			handler: server.CreateUser,
		},
		{
			http:    "GET /user/:username",
			handler: server.GetUser,
		},
		{
			http:    "DELETE /user/:id",
			handler: server.DeleteUser,
		},
		{
			http:    "PUT /user/:id",
			handler: server.UpdateUser,
		},
		{
			http:    "POST /auth",
			handler: server.Login,
		},

		//
		{
			topic:   "send_mail",
			channel: "user",
			handler: server.SendMail,
		},

		//
		{
			event:   "user.created",
			handler: server.Created,
		},

		//
		{
			websocket: "/ws",
			handler:   server.WebSocket,
		},
	}

	/*

		{
			http: "POST /user",
			handler: server.CreateUser,
		},


		{
			method: "POST",
			pattern: "/user",
			handler: server.CreateUser,
		},
		{
			method: "GET",
			pattern: "/user/:username",
			handler: server.GetUser
		},
		{
			method: "DELETE",
			pattern: "/user/:id",
			handler: server.DeleteUser
		},
		{
			method: "PUT",
			pattern: "/user/:id",
			handler: server.UpdateUser
		},
		{
			method: "POST",
			pattern: "/auth"
			handler: server.Login
		}*/

	/*
		{
			topic: "send_mail",
			channel: "user",
			handler: server.SendMail,
		}*/
	/*
		{
			//route: "/es/user.created/",
			event: "user.created",
			handler: server.Created,
		}*/
	/*
		{
			websocket: "/ws",
			handler: server.WebSocket,
		}*/

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

	m := melody.New()
	m.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	g.GET("/websocket", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	go func() {
		for {
			time.Sleep(1 * time.Second)
			m.Broadcast([]byte("Wow"))
		}

	}()

	// The event handlers.
	e.POST("/es/user.created/", "user.created", server.Created)

	return g
}
