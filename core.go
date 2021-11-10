package nb_go

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"time"
)

type CoreCtx struct {
	startTime time.Time
	Provider  IHTTPProvider
	Logger    ILogger

	RootController *HTTPControllerCtx

	Meta  KeyValue
	Setup func()
}

type StartArg struct {
	Host           string
	Port           int
	Path           string
	CORS           *CORSCfg
	ResponseMapper *IResponseMapper
}

type CoreCfg struct {
	Meta *KeyValue
}

func (co *CoreCtx) Boot() {
	co.Setup()
}

func (co *CoreCtx) Start(cfg StartArg) {
	if cfg.ResponseMapper == nil {
		ThrowError(&Err{Message: "ResponseMapper nil!"})
	}

	co.RootController = HTTPController(co.Provider.Router(cfg.Path), co.Logger.NewChild("RootController"), *cfg.ResponseMapper)

	common := CommonController{
		Logger:    co.Logger.NewChild("CommonController"),
		StartTime: time.Now(),
	}

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
}

func Core(config *CoreCfg) *CoreCtx {
	isDebug, _ := StringParser{os.Getenv("DEBUG")}.ToBool()

	validateCoreConfig(config)

	l := Logger("Core")

	if !isDebug {
		gin.SetMode(gin.ReleaseMode)
	}

	p := ExtHTTP()

	c := &CoreCtx{
		Provider:  p,
		startTime: time.Now(),
		Meta:      *config.Meta,
		Logger:    l,
		Setup:     notImplemented("Setup"),
	}

	return c
}
