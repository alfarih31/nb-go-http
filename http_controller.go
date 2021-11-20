package noob

import (
	"errors"
	"fmt"
	"github.com/alfarih31/nb-go-http/app_err"
	"github.com/alfarih31/nb-go-http/keyvalue"
	"github.com/alfarih31/nb-go-http/logger"
	"github.com/alfarih31/nb-go-http/parser"
	"github.com/alfarih31/nb-go-http/tcf"
	"math"
	"net/http"
	"os"
	"runtime"
	"strings"
)

const abortIndex int8 = math.MaxInt8 / 2

type HTTPControllerCtx struct {
	router         *ExtRouter
	Logger         logger.Logger
	ResponseMapper *ResponseMapperCtx
	Debug          bool
}

type ControllerArg struct {
	Logger         logger.Logger
	ResponseMapper *ResponseMapperCtx
}

type HandlerSpec struct {
	method string
	path   string
}

func (h *HTTPControllerCtx) GetSpec(spec string) HandlerSpec {
	if spec == "" {
		apperr.Throw(apperr.New("Handler Spec cannot be empty!"))
	}
	specArr := make([]string, 2, 2)

	for i, s := range strings.Split(spec, " ") {
		specArr[i] = s
		if i == 1 {
			break
		}
	}

	method := specArr[0]
	if method == "" {
		method = "GET"
	}

	path := specArr[1]
	if path == "" {
		path = "/"
	}

	return HandlerSpec{method, path}
}

func (h *HTTPControllerCtx) ToExtHandler(handler HTTPHandler) ExtHandler {
	return func(ec *ExtHandlerCtx) {
		c := WrapExtHandlerCtx(ec)
		tcf.TCFunc(tcf.Func{
			Try: func() {
				res := handler(c)

				if res != nil {
					h.SendSuccess(c, res)
				}
			},
			Catch: func(e interface{}, frames *runtime.Frames) {
				h.SendError(c, e, frames)
			},
		})
	}
}

func (h *HTTPControllerCtx) ToExtHandlers(handlers []HTTPHandler) []ExtHandler {
	extHandlers := make([]ExtHandler, len(handlers), len(handlers))

	for i, handler := range handlers {
		extHandlers[i] = h.ToExtHandler(handler)
	}

	return extHandlers
}

func (h *HTTPControllerCtx) BranchRouter(path string) *ExtRouter {
	h.Logger.Debug(fmt.Sprintf("Branching Controller Router with Path : %s", path), map[string]interface{}{
		"branchFullPath": fmt.Sprintf("%s%s", h.router.fullPath, path),
	})

	return h.router.Branch(path, fmt.Sprintf("%s%s", h.router.fullPath, path))
}

func (h *HTTPControllerCtx) Handle(spec string, handlers ...HTTPHandler) {
	if h.router == nil {
		apperr.Throw(apperr.New("Cannot Set Spec, SetRouter first!"))
	}

	handlerSpec := h.GetSpec(spec)

	switch strings.ToUpper(handlerSpec.method) {
	case "GET":
		h.router.GET(handlerSpec.path, h.ToExtHandlers(handlers)...)
	case "POST":
		h.router.POST(handlerSpec.path, h.ToExtHandlers(handlers)...)
	case "PUT":
		h.router.PUT(handlerSpec.path, h.ToExtHandlers(handlers)...)
	case "DELETE":
		h.router.DELETE(handlerSpec.path, h.ToExtHandlers(handlers)...)
	case "OPTIONS":
		h.router.OPTIONS(handlerSpec.path, h.ToExtHandlers(handlers)...)
	case "HEAD":
		h.router.HEAD(handlerSpec.path, h.ToExtHandlers(handlers)...)
	case "PATCH":
		h.router.PATCH(handlerSpec.path, h.ToExtHandlers(handlers)...)
	case "USE":
		h.router.USE(h.ToExtHandlers(handlers)...)
	default:
		apperr.Throw(apperr.New(fmt.Sprintf("Unknown HTTP Handle Spec: %s", spec)))
	}
}

func (h *HTTPControllerCtx) SendSuccess(c *HandlerCtx, res *Response) {
	r := h.ResponseMapper.GetSuccess()

	r.ComposeTo(res)

	_, rEr := c.response(res.Code, res.Body, res.Header)

	if rEr != nil {
		h.Logger.Error("", map[string]interface{}{"_error": rEr})
	}
}

func (h *HTTPControllerCtx) SendError(c *HandlerCtx, e interface{}, frames *runtime.Frames) {
	r := h.ResponseMapper.GetInternalError()

	// Build Error
	parsedErr := &apperr.AppErr{
		Stack:  apperr.StackTrace(),
		Frames: frames,
	}

	switch er := e.(type) {
	case *apperr.AppErr:
		r = h.ResponseMapper.Get(er.Code, nil)

		// Put to Parsed Error
		parsedErr.Err = er.Err
		parsedErr.Code = er.Code
		parsedErr.Stack = er.Stack

		// If not debug then delete stack from er
		if !h.Debug {
			er.Stack = nil
		}

		if h.Debug {
			r.ComposeBody(keyvalue.KeyValue{"_error": parsedErr.JSON()})
		} else {
			r.ComposeBody(keyvalue.KeyValue{"_error": parsedErr.Error()})
		}
	case Response:
		r = er
		parsedErr.Err = fmt.Errorf("%v", er.Body)
	case *Response:
		r = *er
		parsedErr.Err = fmt.Errorf("%v", er.Body)
	case error:
		parsedErr.Err = er
		if h.Debug {
			r.ComposeBody(keyvalue.KeyValue{"_error": parsedErr.JSON()})
		} else {
			r.ComposeBody(keyvalue.KeyValue{"_error": parsedErr.Error()})
		}
	case string:
		parsedErr.Err = errors.New(er)
		if h.Debug {
			r.ComposeBody(keyvalue.KeyValue{"_error": parsedErr.JSON()})
		} else {
			r.ComposeBody(keyvalue.KeyValue{"_error": parsedErr.Error()})
		}
	default:
		parsedErr.Err = fmt.Errorf("%v", er)
		if h.Debug {
			r.ComposeBody(keyvalue.KeyValue{"_error": parsedErr.JSON()})
		} else {
			r.ComposeBody(keyvalue.KeyValue{"_error": parsedErr.Error()})
		}
	}

	// Stack Error to Context
	c.StackError(parsedErr)

	// If internal error than log error
	if r.Code == http.StatusInternalServerError {
		h.Logger.Error("", map[string]interface{}{"_error": parsedErr.JSON()})
	}

	_, rEr := c.responseError(r.Code, r.Body, r.Header)

	if rEr != nil {
		h.Logger.Error("", map[string]interface{}{"_error": rEr})
	}
}

func (h *HTTPControllerCtx) SetRouter(r *ExtRouter) *HTTPControllerCtx {
	if r == nil {
		apperr.Throw(apperr.New("setRouter r cannot be nil!"))
	}

	h.router = r

	return h
}

func NewController(arg ControllerArg) *HTTPControllerCtx {
	isDebug, _ := parser.String(os.Getenv("DEBUG")).ToBool()

	arg.Logger.Debug("OK", nil)

	h := &HTTPControllerCtx{
		Logger:         arg.Logger,
		ResponseMapper: arg.ResponseMapper,
		Debug:          isDebug,
	}

	return h
}
