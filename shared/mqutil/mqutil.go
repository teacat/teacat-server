package mqutil

import "github.com/gin-gonic/gin"

type Engine struct {
	Gin       *gin.Engine
	Listeners []Listener
}

type Listener struct {
	Method  string
	Path    string
	Channel string
	Topic   string
	Handler gin.HandlerFunc
}

func New(e *gin.Engine) *Engine {
	return &Engine{Gin: e}
}

func (e *Engine) Capture(channel, topic string, handler gin.HandlerFunc) {
	e.Gin.POST("/mq/"+topic, handler)
	e.Listeners = append(e.Listeners, Listener{"POST", "/mq/" + topic, channel, topic, handler})
}
