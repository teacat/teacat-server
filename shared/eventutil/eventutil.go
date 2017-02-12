package eventutil

import "github.com/gin-gonic/gin"

// Engine stores the Gin engine
// and the event listeners.
type Engine struct {
	Gin       *gin.Engine
	Listeners []Listener
}

// Listener listens and processes the specified incoming messages.
type Listener struct {
	Method  string
	Path    string
	Stream  string
	Handler gin.HandlerFunc
}

// New creates the new event listener engine.
func New(e *gin.Engine) *Engine {
	return &Engine{Gin: e}
}

// Capture subscribes to a specified stream and capturing the incoming events.
func (e *Engine) Capture(stream string, handler gin.HandlerFunc) {
	e.Gin.POST("/es/"+stream, handler)
	e.Listeners = append(e.Listeners, Listener{
		Method:  "POST",
		Path:    "/es/" + stream,
		Stream:  stream,
		Handler: handler,
	})
}
