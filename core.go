package nbgohttp

import (
	"context"
	"fmt"
	"github.com/alfarih31/nb-go-http/parser"
	"github.com/gin-gonic/gin"
	"os"
	"time"
)

type CoreCtx struct {
	Context   context.Context
	startTime time.Time
	Provider  *HTTPProviderCtx
	Logger    ILogger

	RootController *HTTPControllerCtx

	Meta  KeyValue
	Setup func()
}

type StartArg struct {
	Host string
	Port int
	Path string
	CORS *CORSCfg
}

type CoreCfg struct {
	Context        context.Context
	Meta           *KeyValue
	ResponseMapper *ResponseMapperCtx
}

func (co *CoreCtx) Boot() {
	co.Setup()
}

func (co *CoreCtx) Start(cfg StartArg) {
	common := CommonController{
		Logger:    co.Logger.NewChild("CommonController"),
		StartTime: time.Now(),
	}

	co.RootController.SetRouter(co.Provider.Router(cfg.Path))

	co.Provider.Engine.NoRoute(co.RootController.ToExtHandlers([]HTTPHandler{common.RequestLogger(), common.HandleNotFound()})...)

	co.RootController.Handle("USE", common.RequestLogger())

	if cfg.CORS != nil {
		if cfg.CORS.Enable {
			co.RootController.Handle("USE", CORS(cfg.CORS))
		}
	}

	co.RootController.Handle("GET /", common.APIStatus(co.Meta))

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

func (c *CoreCtx) WithContext(ctx context.Context) *CoreCtx {
	return Core(&CoreCfg{
		ResponseMapper: c.RootController.ResponseMapper,
		Context:        ctx,
		Meta:           &c.Meta,
	})
}

func Core(config *CoreCfg) *CoreCtx {
	isDebug, _ := parser.String(os.Getenv("DEBUG")).ToBool()

	validateCoreConfig(config)

	l := Logger("Core")

	if !isDebug {
		gin.SetMode(gin.ReleaseMode)
	}

	p := ExtHTTP()

	rc := HTTPController(HTTPControllerArg{
		Logger:         l.NewChild("RootController"),
		ResponseMapper: config.ResponseMapper,
	})

	c := &CoreCtx{
		Context:        config.Context,
		Provider:       p,
		startTime:      time.Now(),
		Meta:           *config.Meta,
		Logger:         l,
		Setup:          notImplemented("Setup"),
		RootController: rc,
	}

	if config.Context != nil {
		c.RootController = c.RootController.WithContext(config.Context)
	}

	return c
}
