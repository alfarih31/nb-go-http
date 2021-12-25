package noob

import (
	_cors "github.com/alfarih31/nb-go-http/cors"
	"net/http"
)

func CORS(config _cors.Cfg) HTTPHandler {
	cors := _cors.New(config)

	return func(c *HandlerCtx) (Response, error) {
		cors.PutCORS(c.Writer)

		if c.Request.Method == http.MethodOptions {
			return nil, nil
		}

		return c.Next()
	}
}
