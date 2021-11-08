package nbgohttp

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HandlerCtx struct {
	ext     *ExtHandlerCtx
	Request *http.Request
	Params  *gin.Params
}

type HTTPHandler func(c *HandlerCtx) *Response

func (c HandlerCtx) Response(status int, body string, headers ResponseHeader) (int, error) {
	if headers != nil {
		for key, head := range headers {
			s, ok := head.(string)
			if ok {
				c.ext.Writer.Header().Set(key, s)
			}
		}
	}

	// Bound status
	if status < 100 || status > 599 {
		status = 500
	}

	c.ext.Writer.WriteHeader(status)
	return c.ext.Writer.WriteString(body)
}

func (c HandlerCtx) Next() {
	c.ext.Next()
}

func WrapExtHandlerCtx(ec *ExtHandlerCtx) *HandlerCtx {
	return &HandlerCtx{
		ext:     ec,
		Request: ec.Request,
		Params:  &ec.Params,
	}
}
