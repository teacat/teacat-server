package mqutil

import "github.com/gin-gonic/gin"

// Engine stores the Gin engine
// and the message queue listeners.
type Engine struct {
	Gin       *gin.Engine
	Listeners []Listener
}

// Listener listens and processes the specified incoming messages.
type Listener struct {
	Method  string
	Path    string
	Channel string
	Topic   string
	Handler gin.HandlerFunc
}

// New creates the new message queue handler engine.
func New(e *gin.Engine) *Engine {
	return &Engine{Gin: e}
}

// Capture subscribes to a channel and capturing the incoming messages with the specified topic.
func (e *Engine) Capture(channel, topic string, handler gin.HandlerFunc) {
	e.Gin.POST("/mq/"+topic, handler)
	e.Listeners = append(e.Listeners, Listener{
		Method:  "POST",
		Path:    "/mq/" + topic,
		Channel: channel,
		Topic:   topic,
		Handler: handler,
	})
}
