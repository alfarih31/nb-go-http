package nbgohttp

import (
	"fmt"
	"time"
)

type CommonController struct {
	Logger    ILogger
	StartTime time.Time
}

func (cc CommonController) APIStatus(m KeyValue) HTTPHandler {
	return func(c *HandlerCtx) *Response {
		u := time.Since(cc.StartTime).String()

		resData := KeyValue{
			"uptime": u,
		}

		resData.Assign(m, true)

		return &Response{
			Body: struct {
				Data interface{} `json:"data"`
			}{
				Data: resData,
			},
		}
	}
}

func (cc CommonController) RequestLogger() HTTPHandler {
	return func(c *HandlerCtx) *Response {
		start := time.Now()

		c.Next()

		cc.Logger.Debug(
			fmt.Sprintf(
				"%s - %s %s %d - %s",
				c.ext.ClientIP(), c.Request.Method, c.Request.URL.Path, c.ext.Writer.Status(), time.Since(start)), map[string]interface{}{
				"clientIp":     c.ext.ClientIP(),
				"method":       c.Request.Method,
				"path":         c.Request.URL.Path,
				"status":       c.ext.Writer.Status(),
				"responseTime": time.Since(start).String(),
			})
		return nil
	}
}

func (cc CommonController) HandleNotFound() HTTPHandler {
	return func(context *HandlerCtx) *Response {
		HTTPError.NotFound.Throw(nil)
		return nil
	}
}

func (cc CommonController) HandleNoMethod() HTTPHandler {
	return func(context *HandlerCtx) *Response {
		HTTPError.NoMethod.Throw(nil)
		return nil
	}
}
