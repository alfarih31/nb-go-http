package noob

import (
	"context"
	"fmt"
	"github.com/alfarih31/nb-go-http/app_err"
	"github.com/alfarih31/nb-go-http/cors"
	"github.com/alfarih31/nb-go-http/keyvalue"
	"github.com/alfarih31/nb-go-http/logger"
	"github.com/alfarih31/nb-go-http/parser"
	"github.com/gin-gonic/gin"
	"net"
	"os"
	"time"
)

type CoreCtx struct {
	startTime time.Time
	Provider  *HTTPProviderCtx
	Logger    logger.Logger

	Meta  keyvalue.KeyValue
	Setup func() // This function will be called when you call the Start of CoreCtx, hence you need to pass the Setup function or the application will be failed to start

	*HTTPControllerCtx
}

type StartArg struct {
	Host       string
	Port       int
	Path       string
	CORS       *cors.Cfg
	Throttling *ThrottlingCfg
	Listener   *net.Listener // Optional use net.Listener if want to start using *net.Listener
}

type CoreCfg struct {
	Context        context.Context
	Meta           *keyvalue.KeyValue
	ResponseMapper *ResponseMapperCtx
}

func (co *CoreCtx) boot() {
	co.Setup()
}

// Start will runt the Core & start serving the application
func (co *CoreCtx) Start(cfg StartArg) {
	common := CommonController{
		Logger:    co.Logger.NewChild("CommonController"),
		StartTime: time.Now(),
	}

	co.SetRouter(co.Provider.Router(cfg.Path))

	co.Provider.Engine.NoRoute(co.ToExtHandlers([]HTTPHandler{common.RequestLogger(), common.HandleNotFound()})...)

	if cfg.Throttling != nil {
		co.Handle("USE", common.Throttling(cfg.Throttling.MaxEventPerSec, cfg.Throttling.MaxEventPerSec))
	}

	co.Handle("USE", common.RequestLogger())

	if cfg.CORS != nil {
		if cfg.CORS.Enable {
			co.Handle("USE", CORS(*cfg.CORS))
		}
	}

	co.boot()

	co.Handle("GET /", common.APIStatus(co.Meta))

	hostInfo := cfg.Host
	if hostInfo == "" {
		hostInfo = "http://localhost"
	}

	baseUrlInfo := fmt.Sprintf("%s:%d", hostInfo, cfg.Port)

	co.Logger.Info(fmt.Sprintf("TimeToBoot = %s Running: BaseUrl = '%s' Path = '%s'", time.Since(co.startTime).String(), baseUrlInfo, cfg.Path), map[string]interface{}{
		"url": fmt.Sprintf("%s%s", baseUrlInfo, cfg.Path),
	})

	var e error
	if cfg.Listener != nil {
		e = co.Provider.Engine.RunListener(*cfg.Listener)
	} else {
		e = co.Provider.Run(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	}

	if e != nil {
		co.Logger.Error("Failed to start, error happened!", map[string]interface{}{"_error": e})
		return
	}

}

func notImplemented(fname string) func() {
	return func() {
		apperr.Throw(apperr.New(fmt.Sprintf("Core.%s Not Implemented", fname)))
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

	p := ExtHTTP()

	rc := NewController(ControllerArg{
		Logger:         l.NewChild("RootController"),
		ResponseMapper: config.ResponseMapper,
	})

	c := &CoreCtx{
		startTime:         time.Now(),
		Provider:          p,
		Meta:              *config.Meta,
		Logger:            l,
		Setup:             notImplemented("Setup"),
		HTTPControllerCtx: rc,
	}

	return c
}
