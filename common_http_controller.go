package noob

import (
	"fmt"
	"github.com/alfarih31/nb-go-http/keyvalue"
	"github.com/alfarih31/nb-go-http/logger"
	"golang.org/x/time/rate"
	"time"
)

type CommonController struct {
	Logger    logger.Logger
	StartTime time.Time
}

type ThrottlingCfg struct {
	MaxEventPerSec int
	MaxBurstSize   int
}

const DefaultMaxBurstSize = 20
const DefaultMaxEventPerSec = 1000

func (cc CommonController) APIStatus(m keyvalue.KeyValue) HTTPHandler {
	return func(c *HandlerCtx) *Response {
		u := time.Since(cc.StartTime).String()

		resData := keyvalue.KeyValue{
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
				c.ClientIP(), c.Request.Method, c.Request.URL.Path, c.Writer.Status(), time.Since(start)), map[string]interface{}{
				"clientIp":     c.ClientIP(),
				"method":       c.Request.Method,
				"path":         c.Request.URL.Path,
				"status":       c.Writer.Status(),
				"responseTime": time.Since(start).String(),
			})
		return nil
	}
}

func (cc CommonController) HandleNotFound() HTTPHandler {
	return func(context *HandlerCtx) *Response {
		HTTPError.NotFound.Throw("")
		return nil
	}
}

func (cc CommonController) Throttling(maxEventsPerSec int, maxBurstSize int) HTTPHandler {
	if maxEventsPerSec == 0 {
		maxEventsPerSec = DefaultMaxEventPerSec
	}

	if maxBurstSize == 0 {
		maxBurstSize = DefaultMaxBurstSize
	}

	limiter := rate.NewLimiter(rate.Limit(maxEventsPerSec), maxBurstSize)

	return func(context *HandlerCtx) *Response {
		if limiter.Allow() {
			context.Next()

			return nil
		}

		HTTPError.TooManyRequest.Throw("")

		return nil
	}
}
