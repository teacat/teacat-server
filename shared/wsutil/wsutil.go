package wsutil

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

type Engine struct {
	Gin *gin.Engine
}

func New(e *gin.Engine) *Engine {
	return &Engine{Gin: e}
}

func (e *Engine) Handle(relativePath string, handler gin.HandlerFunc) {
	m := melody.New()
	m.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	e.Gin.GET("/websocket", func(c *gin.Context) {
		c.Set("websocket", m)
		handler(c)
	})
}
