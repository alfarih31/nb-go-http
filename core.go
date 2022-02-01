package noob

import (
	"fmt"
	"github.com/alfarih31/nb-go-parser"
	"github.com/gin-gonic/gin"
	"net"
	"os"
	"time"
)

type Ctx struct {
	startTime time.Time
	Provider  *HTTPProviderCtx

	Listener net.Listener
	*Router
}

// Start will run the Core & start serving the application
func (co *Ctx) Start() error {
	var (
		e error
	)

	cfg := DefaultCfg

	crs := new(cors)

	// Prepare handlers for no route
	middlewares := []HandlerFunc{handleRequestLogger(Log), crs.HandleCORS, HandleThrottling(), HandleTimeout}

	co.USE(middlewares...)
	// Handle root
	co.GET("/", HandleAPIStatus)

	// handler for not found page
	middlewares = append(middlewares, HandleNotFound)
	co.Provider.Engine.NoRoute(NewHandlerChain(middlewares).compact())

	hostInfo := cfg.Host
	if hostInfo == "" {
		hostInfo = "http://localhost"
	}

	// use listener if listener not nil
	if co.Listener != nil && cfg.UseListener {
		url := fmt.Sprintf("%s%s", co.Listener.Addr().String(), cfg.Path)
		Log.Info(fmt.Sprintf("TimeToBoot = %s Running: Address = '%s'", time.Since(co.startTime).String(), url), map[string]interface{}{
			"address": url,
		})

		e = co.Provider.Engine.RunListener(co.Listener)
	} else {
		baseUrlInfo := fmt.Sprintf("%s:%d", hostInfo, cfg.Port)
		url := fmt.Sprintf("%s%s", baseUrlInfo, cfg.Path)
		Log.Info(fmt.Sprintf("TimeToBoot = %s Running: Url = '%s'", time.Since(co.startTime).String(), url), map[string]interface{}{
			"url": url,
		})

		e = co.Provider.Run(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	}

	return e
}

func notImplemented(fname string) func() error {
	return func() error {
		panic(NewCoreError(fmt.Sprintf("Core.%s not implemented", fname)))

		return nil
	}
}

// New return Core context, used as core of the application
func New(listener ...net.Listener) *Ctx {
	// Load isDebug
	isDebug, _ = parser.String(os.Getenv("DEBUG")).ToBool()

	if !isDebug {
		gin.SetMode(gin.ReleaseMode)
	}

	p := HTTP()

	r := p.Router(DefaultCfg.Path)

	var lis net.Listener
	if len(listener) > 0 {
		lis = listener[0]
	}

	c := &Ctx{
		startTime: time.Now(),
		Provider:  p,
		Listener:  lis,
		Router:    r,
	}

	return c
}
