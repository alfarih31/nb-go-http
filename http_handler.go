package noob

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alfarih31/nb-go-http/utils"
	"github.com/gin-gonic/gin"
	"runtime"
)

const extKeyErrors = "_errors"

var errResponseAlreadyAborted = errors.New("response already aborted")

// Handler is custom type to wrap HandlerFunc, so we can give a name to them
type Handler struct {
	fn   HandlerFunc // fn
	name string      // name of the handler
}

// String return name of the handler. So fmt.Print of a handler will return its name
func (h Handler) String() string {
	if h.name == "" {
		return utils.GetFunctionName(h.fn)
	}

	return h.name
}

// HandlerChain is type for slice of Handler
type HandlerChain []Handler

// Strings return list name of each handler
func (hc HandlerChain) Strings() []string {
	o := make([]string, len(hc))
	for i, h := range hc {
		o[i] = h.String()
	}

	return o
}

func (hc HandlerChain) String() string {
	out := ""
	for i, h := range hc {
		out += fmt.Sprintf("%d. %s\n", i+1, h)
	}

	return out
}

// compact return provider handler chain by compacting all handler to a single gin.HandlerFunc (gin.HandlerFunc)
func (hc HandlerChain) compact() gin.HandlerFunc {
	h := NewHandler(func(c *HandlerCtx) (Response, error) {
		c.setHandlers(hc)

		return c.Next()
	})

	return func(ec *gin.Context) {
		c := WrapHandlerCtx(ec)
		TCFunc(Func{
			Try: func() {
				res, err := h.fn(c)

				// Get response from err
				if err != nil {
					c.SendError(err, GetRuntimeFrames(3))
					return
				}

				if res != nil {
					c.Send(res)
					return
				}
			},
			Catch: func(err interface{}, frames *runtime.Frames) {
				// Send Error
				Log.Warn("error caught! don't panic, use return error instead", map[string]interface{}{"_error": err})

				c.SendError(err, frames)
			},
		})
	}
}

type HandlerCtx struct {
	*gin.Context
	handlerIdx  int
	handlers    HandlerChain
	nextAborted bool
}

type HandlerFunc func(context *HandlerCtx) (Response, error)

func (c *HandlerCtx) response(status HTTPStatusCode, body interface{}, headers map[string]string) error {
	// return if already closed
	if c.nextAborted {
		return nil
	}

	if c.IsAborted() {
		return errResponseAlreadyAborted
	}

	for key, val := range DefaultResponseHeader {
		c.Writer.Header().Set(key, val)
	}

	if headers != nil {
		for key, head := range headers {
			c.Writer.Header().Set(key, head)
		}
	}

	// Bound status
	if status < 100 || status > 599 {
		status = 500
	}

	c.Writer.WriteHeader(int(status))

	j, e := json.Marshal(body)

	if e != nil {
		return e
	}

	_, e = c.Writer.WriteString(string(j))

	// Prevent write to response
	c.Abort()

	return e
}

func (c *HandlerCtx) StackError(e *CoreError) {
	c.Keys[extKeyErrors] = append(c.Keys[extKeyErrors].(Errors), e)
}

func (c *HandlerCtx) Errors() Errors {
	return c.Keys[extKeyErrors].(Errors)
}

func (c *HandlerCtx) setHandlers(handlers HandlerChain) {
	c.handlers = handlers
}

func (c *HandlerCtx) GetNext() HandlerFunc {
	if c.handlerIdx >= len(c.handlers) || c.nextAborted {
		c.Context.Next()

		return nil
	}

	h := c.handlers[c.handlerIdx]
	c.handlerIdx++
	return h.fn
}

func (c *HandlerCtx) Next() (Response, error) {
	h := c.GetNext()

	if h == nil {
		return nil, nil
	}

	res, err := h(c)

	if err != nil {
		c.SendError(err, GetRuntimeFrames(3))
		c.nextAborted = true
		return nil, err
	}

	if res != nil {
		c.Send(res)
		c.nextAborted = true
		return res, nil
	}

	return nil, nil
}

func (c *HandlerCtx) Copy() *HandlerCtx {
	return WrapHandlerCtx(c.Context.Copy())
}

func (c *HandlerCtx) Send(res Response) {
	// Default send success
	r := DefaultSuccessResponse

	// Compose success response to res
	res.Compose(r)

	rEr := c.response(*res.GetCode(), *res.GetBody(), *res.GetHeader())

	if rEr != nil {
		logR.Error("send response error", map[string]interface{}{"_error": rEr})
	}
}

func (c *HandlerCtx) SendError(e interface{}, frames *runtime.Frames) {
	r := DefaultInternalServerErrorResponse

	// Build Error
	parsedErr := &CoreError{
		Stack:  StackTrace(),
		Frames: frames,
	}

	switch er := e.(type) {
	case ResponseError:
		r = er
		parsedErr.Err = er
	case error:
		// Try assertion type to Response
		cR, ok := er.(ResponseError)
		if ok {
			r = cR
		}

		parsedErr.Err = er
	default:
		parsedErr.Err = fmt.Errorf("%v", er)
	}

	// If debug then compose to body
	if isDebug {
		r.ComposeBody(ResponseBody{
			Errors: parsedErr.JSON(),
		})
	} else {
		r.ComposeBody(ResponseBody{
			Errors: parsedErr.Error(),
		})
	}

	// Stack Error to Context
	c.StackError(parsedErr)

	rEr := c.response(*r.GetCode(), *r.GetBody(), *r.GetHeader())

	if rEr != nil {
		logR.Error("send response error", map[string]interface{}{"_error": rEr})
	}
}

func WrapHandlerCtx(ec *gin.Context) *HandlerCtx {
	if ec.Keys == nil {
		ec.Keys = map[string]interface{}{}
	}
	ec.Keys[extKeyErrors] = Errors{}

	return &HandlerCtx{
		Context: ec,
	}
}

func NewHandler(handler HandlerFunc) Handler {
	return Handler{
		name: utils.GetFunctionName(handler),
		fn:   handler,
	}
}

func NewHandlerChain(handlers []HandlerFunc) HandlerChain {
	hs := make([]Handler, len(handlers))
	for i, hnd := range handlers {
		hs[i] = NewHandler(hnd)
	}

	return hs
}
