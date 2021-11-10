package nbgohttp

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HandlerCtx struct {
	Ext     *ExtHandlerCtx
	Request *http.Request
	Params  *gin.Params
}

type HTTPHandler func(context *HandlerCtx) *Response

func (c HandlerCtx) response(status int, body interface{}, headers map[string]string) (int, error) {
	if headers != nil {
		for key, head := range headers {
			c.Ext.Writer.Header().Set(key, head)
		}
	}

	// Bound status
	if status < 100 || status > 599 {
		status = 500
	}

	c.Ext.Writer.WriteHeader(status)

	j, e := json.Marshal(body)

	if e != nil {
		return 0, e
	}

	return c.Ext.Writer.WriteString(string(j))
}

func (c HandlerCtx) Next() {
	c.Ext.Next()
}

func WrapExtHandlerCtx(ec *ExtHandlerCtx) *HandlerCtx {
	return &HandlerCtx{
		Ext:     ec,
		Request: ec.Request,
		Params:  &ec.Params,
	}
}
