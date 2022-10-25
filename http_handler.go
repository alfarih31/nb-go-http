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
const extKeyPrevRes = "_prevRes"
const extKeyPrevErr = "_prevErr"

var errResponseAlreadyAborted = errors.New("response already aborted")

// Handler is custom type to wrap HandlerFunc, so we can give a name to them
type Handler struct {
	fn   HandlerFunc // fn
	name string      // name of the handler
}

// String return name of the handler. So fmt.Print of a handler will return its name. Warning, this features is EXPERIMENTAL!
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
	for _, h := range hc {
		out += fmt.Sprintf("%s,", h)
	}

	return out
}

// compact return provider handler chain by compacting all handler to a single gin.HandlerFunc (gin.HandlerFunc)
func (hc HandlerChain) compact(postHandlers ...HandlerChain) gin.HandlerFunc {
	// h is root of the handler chain
	h := NewHandler(func(c *HandlerCtx) (Response, error) {
		c.setHandlers(hc)

		return c.Next()
	})
	h.name = hc.String()

	var ph *Handler
	if len(postHandlers) > 0 && postHandlers[0] != nil {
		ph = new(Handler)
		*ph = NewHandler(func(c *HandlerCtx) (Response, error) {
			c.setHandlers(postHandlers[0])

			return c.Next()
		})

		ph.name = postHandlers[0].String()
	}

	return func(ec *gin.Context) {
		c := WrapHandlerCtx(ec)
		TCFunc(Func{
			Try: func() {
				// Main handler
				res, err := h.fn(c)

				c.Keys[extKeyPrevRes] = res
				c.Keys[extKeyPrevErr] = err

				if ph != nil {
					res, err = ph.fn(c)
				}

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
				logR.Warn("error caught! don't panic, use return error instead", map[string]interface{}{"_error": err})

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

	// Prevent next handler
	c.nextAborted = true

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
	c.handlerIdx = 0
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

func (c *HandlerCtx) Next() (res Response, err error) {
	h := c.GetNext()

	if h == nil {
		return
	}

	return h(c)
}

func (c *HandlerCtx) Copy() *HandlerCtx {
	return WrapHandlerCtx(c.Context.Copy())
}

func (c *HandlerCtx) Send(res Response) {
	// Default send success
	r := DefaultSuccessResponse.Copy()

	// Compose success response to res
	res.Compose(r)

	rEr := c.response(*res.GetCode(), *res.GetBody(), *res.GetHeader())

	if rEr != nil {
		logR.Error("send response error", map[string]interface{}{"_error": rEr})
	}
}

func (c *HandlerCtx) SendError(e interface{}, frames *runtime.Frames) {
	r := DefaultInternalServerErrorResponse.CopyError()

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

func (c *HandlerCtx) GetPrevResponse() (res Response) {
	if v, exist := c.Keys[extKeyPrevRes]; exist {
		if cv, ok := v.(Response); ok {
			res = cv
		}
	}

	return
}

func (c *HandlerCtx) GetPrevError() (err error) {
	if v, exist := c.Keys[extKeyPrevErr]; exist {
		if cv, ok := v.(error); ok {
			err = cv
		}
	}

	return
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
