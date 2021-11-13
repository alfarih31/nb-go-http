package noob

import (
	"github.com/alfarih31/nb-go-http/cors"
	"net/http"
)

func CORS(config cors.Cfg) HTTPHandler {
	cors := cors.New(config)

	return func(c *HandlerCtx) *Response {
		cors.PutCORS(c.ext.Writer)

		if c.Request.Method == http.MethodOptions {
			return &Response{}
		}

		c.Next()

		return nil
	}
}
