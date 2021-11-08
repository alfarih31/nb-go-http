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
	Config    CoreCfg
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
	CORS CORSCfg
}

type CoreCfg struct {
	Debug bool
	Server
	Meta
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

	co.HTTP.Handle("USE", CORS(co.Config.CORS), common.RequestLogger())
	co.HTTP.Handle("GET /", common.APIStatus(co.Config.Meta))

	co.Boot()

	baseUrl := fmt.Sprintf("%s%d", co.Config.Host, co.Config.Port)

	co.Logger.Debug(fmt.Sprintf("TimeToBoot = %s Running: Url = '%s' Path = '%s'", time.Since(co.startTime).String(), baseUrl, co.Config.Path), nil)
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

func Core(cfg CoreCfg) CoreCtx {
	l := Logger("Core", cfg.Debug)

	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	p := ExtHTTP()
	r := p.Router(cfg.Server.Path)

	h := HTTPController(r, l.NewChild("HTTPController"), ResponseMapper(l.NewChild("ResponseMapper")))

	c := CoreCtx{
		provider:         p,
		startTime:        time.Now(),
		Router:           r,
		HTTP:             h,
		Config:           cfg,
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
