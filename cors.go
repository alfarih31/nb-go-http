package nbgohttp

import (
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
		w.Header().Set(CORSAllowCredentials, BoolParser{cr.Config.AllowCredentials}.ToString())
	}

	if cr.Config.ExposeHeaders != "" {
		w.Header().Set(CORSExposeHeaders, cr.Config.AllowHeaders)
	}
}

func CORS(config *CORSCfg) HTTPHandler {
	cors := TCORS{
		Config: config,
	}

	return func(c *HandlerCtx) *Response {
		cors.PutCORS(c.ext.Writer)

		if c.Request.Method == http.MethodOptions {
			return &Response{}
		}

		c.Next()

		return nil
	}
}
