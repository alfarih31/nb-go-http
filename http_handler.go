package noob

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HandlerCtx struct {
	ext     *ExtHandlerCtx
	Request *http.Request
	Params  *gin.Params
}

type HTTPHandler func(context *HandlerCtx) *Response

func (c HandlerCtx) response(status int, body interface{}, headers map[string]string) (int, error) {
	if headers != nil {
		for key, head := range headers {
			c.ext.Writer.Header().Set(key, head)
		}
	}

	// Bound status
	if status < 100 || status > 599 {
		status = 500
	}

	c.ext.Writer.WriteHeader(status)

	j, e := json.Marshal(body)

	if e != nil {
		return 0, e
	}

	i, e := c.ext.Writer.WriteString(string(j))

	if e != nil {
		return i, e
	}

	c.ext.Abort()

	return 0, nil
}

func (c HandlerCtx) Next() {
	c.ext.Next()
}

func (c HandlerCtx) Query(q string) string {
	return c.ext.Query(q)
}

func WrapExtHandlerCtx(ec *ExtHandlerCtx) *HandlerCtx {
	return &HandlerCtx{
		ext:     ec,
		Request: ec.Request,
		Params:  &ec.Params,
	}
}
