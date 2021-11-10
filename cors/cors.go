package cors

import (
	"github.com/alfarih31/nb-go-http/data"
	http2 "github.com/alfarih31/nb-go-http/http"
	"net/http"
)

const CORSAllowOrigin = "Access-Control-Allow-Origin"
const CORSAllowHeaders = "Access-Control-Allow-Headers"
const CORSAllowMethods = "Access-Control-Allow-Methods"
const CORSAllowCredentials = "Access-Control-Allow-Credentials"
const CORSExposeHeaders = "Access-Control-Expose-Headers"

type CORSCfg struct {
	Enable           bool
	AllowOrigins     string
	AllowMethods     string
	AllowHeaders     string
	AllowCredentials bool
	ExposeHeaders    string
}

type TCORS struct {
	Config *CORSCfg
}

type ICORS interface {
	PutCORS(w http.ResponseWriter)
}

func (cr TCORS) PutCORS(w http.ResponseWriter) {
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
		w.Header().Set(CORSAllowCredentials, data.BoolParser{cr.Config.AllowCredentials}.ToString())
	}

	if cr.Config.ExposeHeaders != "" {
		w.Header().Set(CORSExposeHeaders, cr.Config.AllowHeaders)
	}
}

func CORS(config *CORSCfg) http2.HTTPHandler {
	cors := TCORS{
		Config: config,
	}

	return func(c *http2.HandlerCtx) *http2.Response {
		cors.PutCORS(c.Ext.Writer)

		if c.Request.Method == http.MethodOptions {
			return &http2.Response{}
		}

		c.Next()

		return nil
	}
}
