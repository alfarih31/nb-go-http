package noob

import (
	"github.com/gin-gonic/gin"
	"net"
)

type HTTPProviderCtx struct {
	Engine     *gin.Engine
	rootRouter *Router
}

func (t *HTTPProviderCtx) Router(path string) *Router {
	if path == "/" {
		return t.rootRouter
	}

	return t.rootRouter.Branch(path)
}

func (t *HTTPProviderCtx) preRun() error {
	baseRouter := t.Engine.Group(t.rootRouter.basePath)
	if err := t.rootRouter.boot(baseRouter); err != nil {
		return err
	}
	return nil
}

func (t *HTTPProviderCtx) Run(baseUrl string) error {
	if err := t.preRun(); err != nil {
		return err
	}

	return t.Engine.Run(baseUrl)
}

func (t *HTTPProviderCtx) RunListener(listener net.Listener) error {
	return t.Engine.RunListener(listener)
}

func HTTP() *HTTPProviderCtx {
	h := &HTTPProviderCtx{
		Engine: gin.New(),
		rootRouter: &Router{
			basePath:             "/",
			mapParentMiddlewares: wareCheckers{},
			mapParentPostwares:   wareCheckers{},
		},
	}
	h.Engine.RedirectTrailingSlash = true

	return h
}
