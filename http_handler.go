package noob

import (
	"encoding/json"
	"fmt"
	"github.com/alfarih31/nb-go-http/app_err"
	"github.com/gin-gonic/gin"
	"strings"
)

type Errors []*apperr.AppErr

const extKeyErrors = "_errors"

type HandlerCtx struct {
	*gin.Context
}

type HTTPHandler func(context *HandlerCtx) *Response

func (c *HandlerCtx) response(status int, body interface{}, headers map[string]string) (int, error) {
	if headers != nil {
		for key, head := range headers {
			c.Writer.Header().Set(key, head)
		}
	}

	// Bound status
	if status < 100 || status > 599 {
		status = 500
	}

	c.Writer.WriteHeader(status)

	j, e := json.Marshal(body)

	if e != nil {
		return 0, e
	}

	i, e := c.Writer.WriteString(string(j))

	// Prevent write to response
	c.Abort()

	return i, e
}

func (c *HandlerCtx) StackError(e *apperr.AppErr) {
	c.Keys[extKeyErrors] = append(c.Keys[extKeyErrors].(Errors), e)
}

func (c *HandlerCtx) responseError(status int, e interface{}, headers map[string]string) (int, error) {
	return c.response(status, e, headers)
}

func (c *HandlerCtx) Errors() Errors {
	return c.Keys[extKeyErrors].(Errors)
}

func WrapHandlerCtx(ec *gin.Context) *HandlerCtx {
	ec.Keys = map[string]interface{}{
		extKeyErrors: Errors{},
	}
	return &HandlerCtx{
		Context: ec,
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
