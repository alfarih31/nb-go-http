package http

import (
	"fmt"
	"github.com/alfarih31/nb-go-http/data"
	"github.com/alfarih31/nb-go-http/logger"
	"time"
)

type CommonController struct {
	Logger    logger.ILogger
	StartTime time.Time
}

func (cc CommonController) APIStatus(m data.KeyValue) HTTPHandler {
	return func(c *HandlerCtx) *Response {
		u := time.Since(cc.StartTime).String()

		resData := data.KeyValue{
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
				c.Ext.ClientIP(), c.Request.Method, c.Request.URL.Path, c.Ext.Writer.Status(), time.Since(start)), map[string]interface{}{
				"clientIp":     c.Ext.ClientIP(),
				"method":       c.Request.Method,
				"path":         c.Request.URL.Path,
				"status":       c.Ext.Writer.Status(),
				"responseTime": time.Since(start).String(),
			})
		return nil
	}
}
