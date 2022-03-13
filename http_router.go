package noob

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type wareCheckers map[string]bool

type httpMethod uint

const (
	get httpMethod = iota
	post
	put
	del
	patch
	options
	head
)

// routerHandler is type for Router handler
type routerHandler struct {
	method       httpMethod
	path         string
	handlerChain HandlerChain
}

type Router struct {
	basePath             string
	absPath              string
	handlers             []routerHandler
	middlewares          HandlerChain
	postwares            HandlerChain
	mapParentMiddlewares wareCheckers
	mapParentPostwares   wareCheckers
	branches             []*Router
}

// Handlers return slice to routerHandler
func (e *Router) Handlers() []routerHandler {
	return e.handlers
}

// Branch used for branching router path
func (e *Router) Branch(path string) *Router {
	pm := wareCheckers{}
	for k, v := range e.mapParentMiddlewares {
		pm[k] = v
	}

	pp := wareCheckers{}
	for k, v := range e.mapParentPostwares {
		pp[k] = v
	}
	r := &Router{
		basePath:             path,
		absPath:              fmt.Sprintf("%s%s", e.absPath, path),
		mapParentMiddlewares: pm,
		mapParentPostwares:   pp,
	}

	e.branches = append(e.branches, r)
	return r
}

func (e *Router) GET(path string, handlersFunc ...HandlerFunc) {
	e.handlers = append(e.handlers, routerHandler{
		path:         path,
		method:       get,
		handlerChain: NewHandlerChain(handlersFunc),
	})
}

func (e *Router) POST(path string, handlersFunc ...HandlerFunc) {
	e.handlers = append(e.handlers, routerHandler{
		path:         path,
		method:       post,
		handlerChain: NewHandlerChain(handlersFunc),
	})
}

func (e *Router) PUT(path string, handlersFunc ...HandlerFunc) {
	e.handlers = append(e.handlers, routerHandler{
		path:         path,
		method:       put,
		handlerChain: NewHandlerChain(handlersFunc),
	})
}

func (e *Router) DELETE(path string, handlersFunc ...HandlerFunc) {
	e.handlers = append(e.handlers, routerHandler{
		path:         path,
		method:       del,
		handlerChain: NewHandlerChain(handlersFunc),
	})
}

func (e *Router) PATCH(path string, handlersFunc ...HandlerFunc) {
	e.handlers = append(e.handlers, routerHandler{
		path:         path,
		method:       patch,
		handlerChain: NewHandlerChain(handlersFunc),
	})
}

func (e *Router) OPTIONS(path string, handlersFunc ...HandlerFunc) {
	e.handlers = append(e.handlers, routerHandler{
		path:         path,
		method:       options,
		handlerChain: NewHandlerChain(handlersFunc),
	})
}

func (e *Router) HEAD(path string, handlersFunc ...HandlerFunc) {
	e.handlers = append(e.handlers, routerHandler{
		path:         path,
		method:       head,
		handlerChain: NewHandlerChain(handlersFunc),
	})
}

func (e *Router) USE(handlersFunc ...HandlerFunc) {
	e.middlewares = append(e.middlewares, NewHandlerChain(handlersFunc)...)
}

// POSTUSE assign handler to be executed after other handlers in the groups is executed
func (e *Router) POSTUSE(handlersFunc ...HandlerFunc) {
	e.postwares = append(e.postwares, NewHandlerChain(handlersFunc)...)
}

func (e *Router) boot(parentRouter *gin.RouterGroup) error {
	baseRouter := parentRouter.Group(e.basePath)

	// Filter middlewares to prevent same middlewares invoke twice
	var (
		filteredMiddlewares, filteredPostwares HandlerChain
	)
	for _, m := range e.middlewares {
		// Get middleware name
		mName := m.String()
		if _, exist := e.mapParentMiddlewares[mName]; !exist {
			filteredMiddlewares = append(filteredMiddlewares, m)
			e.mapParentMiddlewares[mName] = true
		}
	}

	for _, m := range e.postwares {
		// Get postware name
		mName := m.String()
		if _, exist := e.mapParentPostwares[mName]; !exist {
			filteredPostwares = append(filteredPostwares, m)
			e.mapParentPostwares[mName] = true
		}
	}

	// put middlewares
	if filteredMiddlewares != nil {
		baseRouter.Use(filteredMiddlewares.compact())
	}

	for _, h := range e.handlers {
		switch h.method {
		case get:
			baseRouter.GET(h.path, h.handlerChain.compact(filteredPostwares))
		case post:
			baseRouter.POST(h.path, h.handlerChain.compact(filteredPostwares))
		case put:
			baseRouter.PUT(h.path, h.handlerChain.compact(filteredPostwares))
		case del:
			baseRouter.DELETE(h.path, h.handlerChain.compact(filteredPostwares))
		case patch:
			baseRouter.PATCH(h.path, h.handlerChain.compact(filteredPostwares))
		case options:
			baseRouter.OPTIONS(h.path, h.handlerChain.compact(filteredPostwares))
		case head:
			baseRouter.HEAD(h.path, h.handlerChain.compact(filteredPostwares))
		}
	}

	for _, b := range e.branches {
		err := b.boot(baseRouter)
		if err != nil {
			return err
		}
	}

	return nil
}
