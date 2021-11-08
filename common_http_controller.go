package nbgohttp

import (
	"fmt"
	"time"
)

type CommonController struct {
	Logger    ILogger
	StartTime time.Time
}

func (cc CommonController) APIStatus(m Meta) HTTPHandler {
	return func(c *HandlerCtx) *Response {
		u := time.Since(cc.StartTime).String()

		data := KeyValue{
			"uptime": u,
		}

		data.Assign(KeyValueFromStruct(m))

		return &Response{
			Body: ResponseBody{
				Data: data,
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
                "clientIp": c.ext.ClientIP(),
                "method": c.Request.Method,
                "path": c.Request.URL.Path,
                "status": c.ext.Writer.Status(),
                "responseTime": time.Since(start),
            })
		return nil
	}
}
