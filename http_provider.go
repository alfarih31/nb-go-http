package noob

import (
	"github.com/gin-gonic/gin"
	"net"
)

type Handler = gin.HandlerFunc
type Router struct {
	path     string
	fullPath string
	router   *gin.RouterGroup
}

func (e *Router) Handlers() []Handler {
	return e.router.Handlers
}

func (e *Router) AppendHandlers(handlers []Handler) {
	finalSize := len(e.router.Handlers) + len(handlers)
	if finalSize >= int(abortIndex) {
		panic("too many handlers")
	}
	mergedHandlers := make([]Handler, finalSize)
	copy(mergedHandlers, e.router.Handlers)
	copy(mergedHandlers[len(e.router.Handlers):], handlers)

	e.router.Handlers = mergedHandlers
}

func (e *Router) Branch(path string, fullPath string) *Router {
	return &Router{
		path:     path,
		fullPath: fullPath,
		router:   e.router.Group(path),
	}
}

func (e *Router) GET(path string, handlers ...Handler) {
	e.router.GET(path, handlers...)
}

func (e *Router) POST(path string, handlers ...Handler) {
	e.router.POST(path, handlers...)
}

func (e *Router) PUT(path string, handlers ...Handler) {
	e.router.PUT(path, handlers...)
}

func (e *Router) DELETE(path string, handlers ...Handler) {
	e.router.DELETE(path, handlers...)
}

func (e *Router) PATCH(path string, handlers ...Handler) {
	e.router.PATCH(path, handlers...)
}

func (e *Router) OPTIONS(path string, handlers ...Handler) {
	e.router.OPTIONS(path, handlers...)
}

func (e *Router) HEAD(path string, handlers ...Handler) {
	e.router.HEAD(path, handlers...)
}

func (e *Router) USE(handlers ...Handler) {
	e.router.Use(handlers...)
}

type HTTPProviderCtx struct {
	Engine *gin.Engine
}

func (t HTTPProviderCtx) Router(path string) *Router {
	return &Router{
		path:     path,
		fullPath: path,
		router:   t.Engine.Group(path),
	}
}

func (t HTTPProviderCtx) Run(url string) error {
	return t.Engine.Run(url)
}

func (t HTTPProviderCtx) RunListener(listener net.Listener) error {
	return t.Engine.RunListener(listener)
}

func HTTP() *HTTPProviderCtx {
	h := &HTTPProviderCtx{
		Engine: gin.New(),
	}

	return h
}
