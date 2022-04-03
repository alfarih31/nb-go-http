package noob

import (
	"context"
	"fmt"
	"github.com/alfarih31/nb-go-keyvalue"
	logger "github.com/alfarih31/nb-go-logger"
	"golang.org/x/time/rate"
	"time"
)

func APIStatus() keyvalue.KeyValue {
	u := time.Since(StartTime).String()

	res := keyvalue.KeyValue{
		"uptime": u,
	}

	res.Assign(DefaultMeta, true)

	return res
}

func HandleAPIStatus(c *HandlerCtx) (Response, error) {
	return NewResponseSuccess(ResponseBody{
		Data: APIStatus(),
	}), nil
}

func handleRequestLogger(logger logger.Logger) HandlerFunc {
	return func(c *HandlerCtx) (Response, error) {
		start := time.Now()

		log := func() {
			latency := time.Since(start)
			logger.Info(
				fmt.Sprintf(
					"%s - %s %s %d - %s",
					c.ClientIP(), c.Request.Method, c.Request.URL.Path, c.Writer.Status(), latency), map[string]interface{}{
					"clientIp": c.ClientIP(),
					"method":   c.Request.Method,
					"path":     c.Request.URL.Path,
					"status":   c.Writer.Status(),
					"latency":  latency.String(),
				})
		}

		defer log()

		return c.Next()
	}
}

func HandleNotFound(context *HandlerCtx) (Response, error) {
	// Don't handle if http version > 1
	if context.Request.ProtoMajor > 1 {
		return context.Next()
	}

	return nil, DefaultNotFoundErrorResponse
}

func HandleTimeout(c *HandlerCtx) (Response, error) {
	if DefaultCfg.RequestTimeout > 0 {
		timeoutCtx, cancel := context.WithTimeout(c.Request.Context(), DefaultCfg.RequestTimeout)
		defer cancel()

		resChan := make(chan Response)
		errChan := make(chan error)

		go func() {
			res, err := c.Next()

			if err != nil {
				errChan <- err
			}

			resChan <- res
		}()

		select {
		case res := <-resChan:
			return res, nil
		case err := <-errChan:
			return nil, err
		case <-timeoutCtx.Done():
			r := DefaultRequestTimeoutErrorResponse
			if err := c.response(*r.GetCode(), *r.GetBody(), *r.GetHeader()); err != nil {
				log.Error(err)
			}
		}

		return nil, nil
	}

	return c.Next()
}

func HandleThrottling() HandlerFunc {
	cfg := DefaultThrottlingCfg

	if !cfg.Enable {
		return func(context *HandlerCtx) (Response, error) {
			return context.Next()
		}
	}

	limiter := rate.NewLimiter(rate.Limit(cfg.MaxEventPerSec), cfg.MaxBurstSize)

	return func(context *HandlerCtx) (Response, error) {
		if limiter.Allow() {
			return context.Next()
		}
		return nil, DefaultTooManyRequestsErrorResponse
	}
}
