package cors

import (
	"github.com/alfarih31/nb-go-parser"
	"net/http"
)

const CORSAllowOrigin = "Access-Control-Allow-Origin"
const CORSAllowHeaders = "Access-Control-Allow-Headers"
const CORSAllowMethods = "Access-Control-Allow-Methods"
const CORSAllowCredentials = "Access-Control-Allow-Credentials"
const CORSExposeHeaders = "Access-Control-Expose-Headers"

type Cfg struct {
	Enable           bool
	AllowOrigins     string
	AllowMethods     string
	AllowHeaders     string
	AllowCredentials bool
	ExposeHeaders    string
}

type cors struct {
	Config Cfg
}

type CORS interface {
	PutCORS(w http.ResponseWriter)
}

func (cr cors) PutCORS(w http.ResponseWriter) {
	if cr.Config.AllowOrigins != "" {
		w.Header().Set(CORSAllowOrigin, cr.Config.AllowOrigins)
	}

	if cr.Config.AllowHeaders != "" {
		w.Header().Set(CORSAllowMethods, cr.Config.AllowMethods)
	}

	if cr.Config.AllowHeaders != "" {
		w.Header().Set(CORSAllowHeaders, cr.Config.AllowHeaders)
	}

	if cr.Config.AllowCredentials == true {
		w.Header().Set(CORSAllowCredentials, parser.Bool(cr.Config.AllowCredentials).ToString())
	}

	if cr.Config.ExposeHeaders != "" {
		w.Header().Set(CORSExposeHeaders, cr.Config.AllowHeaders)
	}
}

func New(config Cfg) cors {
	return cors{
		Config: config,
	}
}
