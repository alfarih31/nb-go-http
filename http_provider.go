package nbgohttp

import (
	"github.com/gin-gonic/gin"
)

type ExtHandlerCtx = gin.Context
type ExtHandler = gin.HandlerFunc
type ExtRouter struct {
	router *gin.RouterGroup
}

func (e *ExtRouter) Handlers() []ExtHandler {
	return e.router.Handlers
}

func (e *ExtRouter) AppendHandlers(handlers []ExtHandler) {
	finalSize := len(e.router.Handlers) + len(handlers)
	if finalSize >= int(abortIndex) {
		panic("too many handlers")
	}
	mergedHandlers := make([]ExtHandler, finalSize)
	copy(mergedHandlers, e.router.Handlers)
	copy(mergedHandlers[len(e.router.Handlers):], handlers)

	e.router.Handlers = mergedHandlers
}

func (e *ExtRouter) Branch(path string) *ExtRouter {
	return &ExtRouter{
		router: e.router.Group(path),
	}
}

func (e *ExtRouter) GET(path string, handlers ...ExtHandler) {
	e.router.GET(path, handlers...)
}

func (e *ExtRouter) POST(path string, handlers ...ExtHandler) {
	e.router.POST(path, handlers...)
}

func (e *ExtRouter) PUT(path string, handlers ...ExtHandler) {
	e.router.PUT(path, handlers...)
}

func (e *ExtRouter) DELETE(path string, handlers ...ExtHandler) {
	e.router.DELETE(path, handlers...)
}

func (e *ExtRouter) PATCH(path string, handlers ...ExtHandler) {
	e.router.PATCH(path, handlers...)
}

func (e *ExtRouter) OPTIONS(path string, handlers ...ExtHandler) {
	e.router.OPTIONS(path, handlers...)
}

func (e *ExtRouter) HEAD(path string, handlers ...ExtHandler) {
	e.router.HEAD(path, handlers...)
}

func (e *ExtRouter) USE(handlers ...ExtHandler) {
	e.router.Use(handlers...)
}

type IHTTPProvider interface {
	Router(path string) *ExtRouter
	Run(url string) error
}

type THTTPProvider struct {
	Engine *gin.Engine
}

func (t THTTPProvider) Router(path string) *ExtRouter {
	return &ExtRouter{
		router: t.Engine.Group(path),
	}
}

func (t THTTPProvider) Run(url string) error {
	return t.Engine.Run(url)
}

func ExtHTTP() IHTTPProvider {
	h := THTTPProvider{
		Engine: gin.New(),
	}

	return h
}
