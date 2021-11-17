package noob

import (
	"encoding/json"
	"fmt"
	"github.com/alfarih31/nb-go-http/app_err"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type Errors []*apperr.AppErr

type HandlerCtx struct {
	ext      *ExtHandlerCtx
	Request  *http.Request
	Response gin.ResponseWriter
	Params   *gin.Params
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

func (c *HandlerCtx) StackError(e *apperr.AppErr) {
	errs, exist := c.ext.Keys["errors"]
	if !exist {
		errs = Errors{}
	}

	errs = append(errs.(Errors), e)
	c.ext.Keys["errors"] = errs
}

func (c *HandlerCtx) responseError(status int, e interface{}, headers map[string]string) (int, error) {
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

func (c *HandlerCtx) Errors() Errors {
	return c.ext.Keys["errors"].(Errors)
}

func WrapExtHandlerCtx(ec *ExtHandlerCtx) *HandlerCtx {
	ec.Keys = map[string]interface{}{}
	return &HandlerCtx{
		ext:      ec,
		Request:  ec.Request,
		Response: ec.Writer,
		Params:   &ec.Params,
	}
}

func (e Errors) MarshalJSON() ([]byte, error) {
	jsonData := make([]interface{}, len(e))
	for i, er := range e {
		jsonData[i] = er.JSON()
	}

	return json.Marshal(jsonData)
}

func (e Errors) String() string {
	if len(e) == 0 {
		return ""
	}
	var buffer strings.Builder
	for i, msg := range e {
		fmt.Fprintf(&buffer, "Error #%02d: %v\n", i+1, msg)
	}
	return buffer.String()
}
