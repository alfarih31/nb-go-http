package noob

import (
	"github.com/alfarih31/nb-go-parser"
	"net/http"
	"strconv"
	"time"
)

const CORSAllowOrigin = "Access-Control-Allow-Origin"
const CORSAllowHeaders = "Access-Control-Allow-Headers"
const CORSAllowMethods = "Access-Control-Allow-Methods"
const CORSAllowCredentials = "Access-Control-Allow-Credentials"
const CORSExposeHeaders = "Access-Control-Expose-Headers"
const CORSMaxAge = "Access-Control-Max-Age"

type cors struct {
}

func (c cors) validateOrigins(origin string) bool {
	for _, o := range DefaultCORSCfg.AllowOrigins {
		if o == "*" {
			return true
		}

		if o == origin {
			return true
		}
	}

	return false
}

func (c cors) HandleCORS(ctx *HandlerCtx) (Response, error) {
	if !DefaultCORSCfg.Enable {
		return ctx.Next()
	}

	origin := ctx.GetHeader("Origin")
	if origin == "" {
		return ctx.Next()
	}

	cfg := DefaultCORSCfg

	if cfg.AllowOrigins == nil {
		ctx.Writer.Header().Set(CORSAllowOrigin, origin)
	} else {

		if !c.validateOrigins(origin) {
			return DefaultForbiddenErrorResponse, nil
		}

		ctx.Writer.Header().Set(CORSAllowOrigin, origin)
	}

	if ctx.Request.Method == http.MethodOptions {
		c.handlePreflightRequest(ctx)
		return DefaultSuccessNoContentResponse, nil
	}

	c.handleNormalRequest(ctx)

	return ctx.Next()
}

func (c cors) handlePreflightRequest(ctx *HandlerCtx) {
	c.applyPreflightHeaders(ctx.Writer)
}

func (c cors) handleNormalRequest(ctx *HandlerCtx) {
	c.applyNormalHeaders(ctx.Writer)
}

func (c cors) applyNormalHeaders(w http.ResponseWriter) {
	cfg := DefaultCORSCfg

	if cfg.AllowCredentials {
		w.Header().Set(CORSAllowCredentials, parser.Bool(cfg.AllowCredentials).ToString())
	}

	if cfg.ExposeHeaders != "" {
		w.Header().Set(CORSExposeHeaders, cfg.ExposeHeaders)
	}

	w.Header().Set("Vary", "Origin")
}

func (c cors) applyPreflightHeaders(w http.ResponseWriter) {
	cfg := DefaultCORSCfg

	if cfg.AllowMethods != "" {
		w.Header().Set(CORSAllowMethods, cfg.AllowMethods)
	}

	if cfg.AllowHeaders != "" {
		w.Header().Set(CORSAllowHeaders, cfg.AllowHeaders)
	}

	if cfg.AllowCredentials {
		w.Header().Set(CORSAllowCredentials, parser.Bool(cfg.AllowCredentials).ToString())
	}

	if cfg.ExposeHeaders != "" {
		w.Header().Set(CORSExposeHeaders, cfg.ExposeHeaders)
	}

	if cfg.MaxAge > time.Duration(0) {
		w.Header().Set(CORSMaxAge, strconv.FormatInt(int64(cfg.MaxAge.Seconds()), 10))
	}

	w.Header().Add("Vary", "Origin")
	w.Header().Add("Vary", "Access-Control-Request-Method")
	w.Header().Add("Vary", "Access-Control-Request-Headers")
}
