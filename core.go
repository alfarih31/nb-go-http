package noob

import (
	"context"
	"fmt"
	"github.com/alfarih31/nb-go-http/keyvalue"
	"github.com/alfarih31/nb-go-http/logger"
	"github.com/alfarih31/nb-go-http/parser"
	"github.com/gin-gonic/gin"
	"os"
	"time"
)

type coreCtx struct {
	startTime time.Time
	Provider  *HTTPProviderCtx
	Logger    logger.ILogger

	Meta  keyvalue.KeyValue
	Setup func()

	*httpControllerCtx
}

type StartArg struct {
	Host string
	Port int
	Path string
	CORS *CORSCfg
}

type CoreCfg struct {
	Context        context.Context
	Meta           *keyvalue.KeyValue
	ResponseMapper *responseMapperCtx
}

func (co *coreCtx) Boot() {
	co.Setup()
}

func (co *coreCtx) Start(cfg StartArg) {
	common := CommonController{
		Logger:    co.Logger.NewChild("CommonController"),
		StartTime: time.Now(),
	}

	co.SetRouter(co.Provider.Router(cfg.Path))

	co.Provider.Engine.NoRoute(co.ToExtHandlers([]HTTPHandler{common.RequestLogger(), common.HandleNotFound()})...)

	co.Handle("USE", common.RequestLogger())

	if cfg.CORS != nil {
		if cfg.CORS.Enable {
			co.Handle("USE", CORS(cfg.CORS))
		}
	}

	co.Handle("GET /", common.APIStatus(co.Meta))

	co.Boot()

	if cfg.Host == "" {
		cfg.Host = ":"
	}

	hostInfo := cfg.Host
	if hostInfo == ":" {
		hostInfo = "http://localhost"
	}

	baseUrlInfo := fmt.Sprintf("%s:%d", hostInfo, cfg.Port)

	co.Logger.Info(fmt.Sprintf("TimeToBoot = %s Running: BaseUrl = '%s' Path = '%s'", time.Since(co.startTime).String(), baseUrlInfo, cfg.Path), map[string]interface{}{
		"url": fmt.Sprintf("%s%s", baseUrlInfo, cfg.Path),
	})
	e := co.Provider.Run(fmt.Sprintf("%s%d", cfg.Host, cfg.Port))

	if e != nil {
		co.Logger.Error("Failed to start, error happened!", map[string]interface{}{"_error": e})
		return
	}

}

func notImplemented(fname string) func() {
	return func() {
		ThrowError(&Err{
			Message: fmt.Sprintf("Core.%s Not Implemented", fname),
		})
	}
}

func validateCoreConfig(config *CoreCfg) {
	if config == nil {
		ThrowError(&Err{Message: "Core config cannot be nil"})
	}

	if config.Meta == nil {
		ThrowError(&Err{Message: "Core config.Meta cannot be nil"})
	}

	if config.ResponseMapper == nil {
		ThrowError(&Err{Message: "Core config.ResponseMapper cannot be nil"})
	}
}

func New(config *CoreCfg) *coreCtx {
	isDebug, _ := parser.String(os.Getenv("DEBUG")).ToBool()

	validateCoreConfig(config)

	l := logger.Logger("Core")

	if !isDebug {
		gin.SetMode(gin.ReleaseMode)
	}

	p := ExtHTTP()

	rc := NewController(ControllerArg{
		Logger:         l.NewChild("RootController"),
		ResponseMapper: config.ResponseMapper,
	})

	c := &coreCtx{
		Provider:          p,
		Meta:              *config.Meta,
		Logger:            l,
		Setup:             notImplemented("Setup"),
		httpControllerCtx: rc,
	}

	return c
}
