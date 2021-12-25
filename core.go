package noob

import (
	"context"
	"fmt"
	"github.com/alfarih31/nb-go-http/app_err"
	"github.com/alfarih31/nb-go-http/cors"
	"github.com/alfarih31/nb-go-http/logger"
	"github.com/alfarih31/nb-go-keyvalue"
	"github.com/alfarih31/nb-go-parser"
	"github.com/gin-gonic/gin"
	"net"
	"os"
	"time"
)

type CoreCtx struct {
	startTime time.Time
	Provider  *HTTPProviderCtx
	Logger    logger.Logger

	Meta     keyvalue.KeyValue
	Setup    func() error
	Listener net.Listener

	HTTPController
}

type StartArg struct {
	Host        string
	Port        int
	Path        string
	CORS        *cors.Cfg
	Throttling  *ThrottlingCfg
	UseListener bool
}

type CoreCfg struct {
	Context        context.Context
	Meta           keyvalue.KeyValue
	ResponseMapper *ResponseMapperCtx
	Listener       net.Listener // Optional use net.Listener if want to start using *net.Listener
}

func (co *CoreCtx) boot() error {
	if err := co.Setup(); err != nil {
		return err
	}

	return nil
}

// Start will runt the Core & start serving the application
func (co *CoreCtx) Start(cfg StartArg) {
	var (
		e error
	)

	common := CommonController{
		Logger:    co.Logger.NewChild("CommonController"),
		StartTime: time.Now(),
	}

	co.SetRouter(co.Provider.Router(cfg.Path))

	co.Provider.Engine.NoRoute(co.chainHandlers([]HTTPHandler{common.RequestLogger(), common.HandleNotFound()}))

	if cfg.Throttling != nil {
		co.Handle("USE", common.Throttling(cfg.Throttling.MaxEventPerSec, cfg.Throttling.MaxEventPerSec))
	}

	co.Handle("USE", common.RequestLogger())

	if cfg.CORS != nil {
		if cfg.CORS.Enable {
			co.Handle("USE", CORS(*cfg.CORS))
		}
	}

	e = co.boot()
	if e != nil {
		apperr.Throw(apperr.New("App failed to boot", e))
	}

	co.Handle("GET /", common.APIStatus(co.Meta))

	hostInfo := cfg.Host
	if hostInfo == "" {
		hostInfo = "http://localhost"
	}

	// If use listener then start using listener
	if co.Listener != nil && cfg.UseListener {
		url := fmt.Sprintf("%s%s", co.Listener.Addr().String(), cfg.Path)
		co.Logger.Info(fmt.Sprintf("TimeToBoot = %s Running: Address = '%s'", time.Since(co.startTime).String(), url), map[string]interface{}{
			"address": url,
		})

		e = co.Provider.Engine.RunListener(co.Listener)
	} else {
		baseUrlInfo := fmt.Sprintf("%s:%d", hostInfo, cfg.Port)
		url := fmt.Sprintf("%s%s", baseUrlInfo, cfg.Path)
		co.Logger.Info(fmt.Sprintf("TimeToBoot = %s Running: Url = '%s'", time.Since(co.startTime).String(), url), map[string]interface{}{
			"url": url,
		})

		e = co.Provider.Run(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	}

	if e != nil {
		co.Logger.Error("Failed to start, error happened!", map[string]interface{}{"_error": e})
		return
	}

}

func notImplemented(fname string) func() error {
	return func() error {
		apperr.Throw(apperr.New(fmt.Sprintf("Core.%s Not Implemented", fname)))

		return nil
	}
}

func validateCoreConfig(config *CoreCfg) {
	if config == nil {
		apperr.Throw(apperr.New("Core config cannot be nil"))
	}

	if config.Meta == nil {
		apperr.Throw(apperr.New("Core config.Meta cannot be nil"))
	}

	if config.ResponseMapper == nil {
		apperr.Throw(apperr.New("Core config.ResponseMapper cannot be nil"))
	}
}

// New return Core context, used as core of the application
func New(config *CoreCfg) *CoreCtx {
	isDebug, _ := parser.String(os.Getenv("DEBUG")).ToBool()

	validateCoreConfig(config)

	l := logger.New("Core")

	if !isDebug {
		gin.SetMode(gin.ReleaseMode)
	}

	p := HTTP()

	rc := NewHTTPController(ControllerArg{
		Logger:         l.NewChild("RootController"),
		ResponseMapper: config.ResponseMapper,
	})

	c := &CoreCtx{
		startTime:      time.Now(),
		Provider:       p,
		Meta:           config.Meta,
		Logger:         l,
		Setup:          notImplemented("Setup"),
		HTTPController: rc,
		Listener:       config.Listener,
	}

	return c
}
