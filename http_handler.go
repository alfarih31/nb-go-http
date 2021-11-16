package noob

import (
	"encoding/json"
	"fmt"
	"github.com/alfarih31/nb-go-http/app_err"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HandlerCtx struct {
	ext      *ExtHandlerCtx
	Request  *http.Request
	Response gin.ResponseWriter
	Params   *gin.Params
	Errors   *[]*gin.Error
}

type HTTPHandler func(context *HandlerCtx) *Response

func (c *HandlerCtx) response(status int, body interface{}, headers map[string]string) (int, error) {
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

	return c.ext.Writer.WriteString(string(j))
}

func (c *HandlerCtx) responseError(status int, e interface{}, headers map[string]string) (int, error) {
	switch er := e.(type) {
	case error:
		_ = c.ext.Error(er)
	case apperr.AppErr:
		_ = c.ext.Error(er)
	case *apperr.AppErr:
		_ = c.ext.Error(er)
	default:
		_ = c.ext.Error(fmt.Errorf("%v", er))
	}

	c.Errors = (*[]*gin.Error)(&c.ext.Errors)

	i, err := c.response(status, e, headers)
	if err != nil {
		return i, err
	}

	c.ext.Abort()

	return i, nil
}

func (c *HandlerCtx) Next() {
	c.ext.Next()
}

func (c *HandlerCtx) Query(q string) string {
	return c.ext.Query(q)
}

func WrapExtHandlerCtx(ec *ExtHandlerCtx) *HandlerCtx {
	return &HandlerCtx{
		ext:      ec,
		Request:  ec.Request,
		Response: ec.Writer,
		Params:   &ec.Params,
		Errors:   (*[]*gin.Error)(&ec.Errors),
	}
}
