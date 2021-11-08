package nbgohttp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

type CoreCtx struct {
	startTime time.Time
	provider  IHTTPProvider
	Router    *ExtRouter
	Config    *CoreCfg
	Logger    ILogger

	HTTP           HTTPControllerCtx
	InitComponents func()

	InitDatasource   func()
	InitRepositories func()
	InitServices     func()
	InitControllers  func()
	Setup            func()
}

type Server struct {
	Host string
	Port int
	Path string
	CORS *CORSCfg
}

type CoreCfg struct {
	Debug  bool
	Server *Server
	Meta   *Meta
}

func (co CoreCtx) Boot() {
	co.Setup()
	co.InitComponents()
	co.InitDatasource()
	co.InitRepositories()
	co.InitServices()
	co.InitControllers()
}

func (co CoreCtx) Start() {
	common := CommonController{
		Logger:    co.Logger.NewChild("CommonController"),
		StartTime: time.Now(),
	}

	if co.Config.Server.CORS != nil {
		if co.Config.Server.CORS.Enable {
			co.HTTP.Handle("USE", CORS(co.Config.Server.CORS))
		}
	}

	co.HTTP.Handle("USE", common.RequestLogger())
	co.HTTP.Handle("GET /", common.APIStatus(*co.Config.Meta))

	co.Boot()

	baseUrl := fmt.Sprintf("%s%d", co.Config.Server.Host, co.Config.Server.Port)

	co.Logger.Debug(fmt.Sprintf("TimeToBoot = %s Running: Url = '%s' Path = '%s'", time.Since(co.startTime).String(), baseUrl, co.Config.Server.Path), nil)
	e := co.provider.Run(baseUrl)

	if e != nil {
		co.Logger.Error("Failed to start, error happened!", map[string]interface{}{"_error": e})
		return
	}

}

func notImplemented(fname string) func() {
	return func() {
		ThrowError(Err{
			Message: fmt.Sprintf("%s Not Implemented", fname),
		})
	}
}

func validateCoreConfig(config *CoreCfg) {
	if config == nil {
		ThrowError(Err{Message: "Core config cannot be nil"})
	}

	if config.Server == nil {
		ThrowError(Err{Message: "Core config.Server cannot be nil"})
	}

	if config.Meta == nil {
		ThrowError(Err{Message: "Core config.Meta cannot be nil"})
	}
}

func Core(config *CoreCfg) CoreCtx {
	validateCoreConfig(config)

	l := Logger("Core", config.Debug)

	if !config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	p := ExtHTTP()
	r := p.Router(config.Server.Path)

	h := HTTPController(r, l.NewChild("HTTPController"), ResponseMapper(l.NewChild("ResponseMapper")))

	c := CoreCtx{
		provider:         p,
		startTime:        time.Now(),
		Router:           r,
		HTTP:             h,
		Config:           config,
		Logger:           l,
		Setup:            notImplemented("Setup"),
		InitComponents:   notImplemented("Init Components"),
		InitDatasource:   notImplemented("Init Datasource"),
		InitRepositories: notImplemented("Init Repositories"),
		InitServices:     notImplemented("Init Services"),
		InitControllers:  notImplemented("Init Controllers"),
	}

	return c
}
