package eventutil

import "github.com/gin-gonic/gin"

type Engine struct {
	Listeners []Listener
}

type HandlerFunc func(*gin.Context)

type Listener struct {
	Method  string
	Path    string
	Stream  string
	Handler HandlerFunc
}

func New() *Engine {
	return &Engine{}
}

func (e *Engine) POST(relativePath string, stream string, handler HandlerFunc) {
	e.Listeners = append(e.Listeners, Listener{"POST", relativePath, stream, handler})
}

func (e *Engine) GET(relativePath string, stream string, handler HandlerFunc) {
	e.Listeners = append(e.Listeners, Listener{"GET", relativePath, stream, handler})
}

func (e *Engine) DELETE(relativePath string, stream string, handler HandlerFunc) {
	e.Listeners = append(e.Listeners, Listener{"DELETE", relativePath, stream, handler})
}

func (e *Engine) PUT(relativePath string, stream string, handler HandlerFunc) {
	e.Listeners = append(e.Listeners, Listener{"PUT", relativePath, stream, handler})
}

func (e *Engine) PATCH(relativePath string, stream string, handler HandlerFunc) {
	e.Listeners = append(e.Listeners, Listener{"PATCH", relativePath, stream, handler})
}
