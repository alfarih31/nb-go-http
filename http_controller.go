package noob

import (
	"fmt"
	"github.com/alfarih31/nb-go-http/app_err"
	"github.com/alfarih31/nb-go-http/logger"
	"github.com/alfarih31/nb-go-http/tcf"
	"github.com/alfarih31/nb-go-keyvalue"
	"github.com/alfarih31/nb-go-parser"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"os"
	"runtime"
	"strings"
)

const abortIndex int8 = math.MaxInt8 / 2

type HTTPControllerCtx struct {
	logger.Logger
	router         *Router
	ResponseMapper *ResponseMapperCtx
	IsDebug        bool
}

type HTTPController interface {
	logger.Logger
	SendError(c *HandlerCtx, e interface{}, frames *runtime.Frames)
	Send(c *HandlerCtx, res Response)
	GetSpec(spec string) HandlerSpec
	AdaptHandler(handler HTTPHandler) Handler
	BranchRouter(path string) *Router
	Handle(spec string, handlers ...HTTPHandler)
	SetRouter(r *Router) *HTTPControllerCtx
	chainHandlers(handlers []HTTPHandler) Handler
}

type ControllerArg struct {
	Logger         logger.Logger
	ResponseMapper *ResponseMapperCtx
	Router         *Router
}

type HandlerSpec struct {
	method string
	path   string
}

func (h *HTTPControllerCtx) GetSpec(spec string) HandlerSpec {
	specArr := strings.SplitN(spec, " ", 2)
	if len(specArr) < 1 {
		apperr.Throw(apperr.New("Handler spec wrong format, expected: {METHOD} {PATH}, example: \"GET /foo\""))
	}

	method := specArr[0]

	if strings.ToUpper(method) == "USE" {
		return HandlerSpec{method, "/"}
	}

	path := specArr[1]

	return HandlerSpec{method, path}
}

func (h *HTTPControllerCtx) AdaptHandler(handler HTTPHandler) Handler {
	return func(ec *gin.Context) {
		c := WrapHandlerCtx(ec)
		tcf.TCFunc(tcf.Func{
			Try: func() {
				res, err := handler(c)

				if err != nil {
					h.SendError(c, err, apperr.GetRuntimeFrames(3))
					return
				}

				if res == nil {
					return
				}

				h.Send(c, res)
			},
			Catch: func(e interface{}, frames *runtime.Frames) {
				// Send Error
				h.SendError(c, e, frames)
			},
		})
	}
}

func (h *HTTPControllerCtx) BranchRouter(path string) *Router {
	h.Logger.Debug(fmt.Sprintf("Branching Controller Router with Path : %s", path), map[string]interface{}{
		"branchFullPath": fmt.Sprintf("%s%s", h.router.fullPath, path),
	})

	return h.router.Branch(path, fmt.Sprintf("%s%s", h.router.fullPath, path))
}

func (h *HTTPControllerCtx) chainHandlers(handlers []HTTPHandler) Handler {
	return h.AdaptHandler(func(context *HandlerCtx) (Response, error) {
		context.handlers = handlers

		res, err := handlers[0](context)

		return res, err
	})
}

func (h *HTTPControllerCtx) Handle(spec string, handlers ...HTTPHandler) {
	if h.router == nil {
		apperr.Throw(apperr.New("Cannot Set Spec, SetRouter first!"))
	}

	handlerSpec := h.GetSpec(spec)

	switch strings.ToUpper(handlerSpec.method) {
	case "GET":
		h.router.GET(handlerSpec.path, h.chainHandlers(handlers))
	case "POST":
		h.router.POST(handlerSpec.path, h.chainHandlers(handlers))
	case "PUT":
		h.router.PUT(handlerSpec.path, h.chainHandlers(handlers))
	case "DELETE":
		h.router.DELETE(handlerSpec.path, h.chainHandlers(handlers))
	case "OPTIONS":
		h.router.OPTIONS(handlerSpec.path, h.chainHandlers(handlers))
	case "HEAD":
		h.router.HEAD(handlerSpec.path, h.chainHandlers(handlers))
	case "PATCH":
		h.router.PATCH(handlerSpec.path, h.chainHandlers(handlers))
	case "USE":
		h.router.USE(h.chainHandlers(handlers))
	default:
		apperr.Throw(apperr.New(fmt.Sprintf("Unknown HTTP Handle Spec: %s", spec)))
	}
}

func (h *HTTPControllerCtx) Send(c *HandlerCtx, res Response) {
	// Default send success
	r := h.ResponseMapper.GetSuccess()

	// Compose success response to res
	res.Compose(r)

	// Set default header
	h.ResponseMapper.AssignDefaultHeader(res)

	_, rEr := c.response(res.GetCode(), res.GetBody(), res.GetHeader())

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
		r = h.ResponseMapper.Get(er.Code)

		// Put to Parsed Error
		parsedErr.Err = er.Err
		parsedErr.Code = er.Code
		parsedErr.Stack = er.Stack

		// If not debug then delete stack from er
		if !h.IsDebug {
			er.Stack = nil
		}
	case error:
		// Try assertion type to Response
		cR, ok := er.(Response)
		if ok {
			r = cR
		}

		parsedErr.Err = er
	default:
		parsedErr.Err = fmt.Errorf("%v", er)
	}

	// If debug then compose to body
	if h.IsDebug {
		r.ComposeBody(keyvalue.KeyValue{"_error": parsedErr.JSON()})
	} else {
		r.ComposeBody(keyvalue.KeyValue{"_error": parsedErr.Error()})
	}

	// Stack Error to Context
	c.StackError(parsedErr)

	// If internal error than log error
	if r.GetCode() == http.StatusInternalServerError {
		h.Logger.Error("", map[string]interface{}{"_error": parsedErr.JSON()})
	}

	// Set default header
	h.ResponseMapper.AssignDefaultHeader(r)

	_, rEr := c.responseError(r.GetCode(), r.GetBody(), r.GetHeader())

	if rEr != nil {
		h.Logger.Error("", map[string]interface{}{"_error": rEr})
	}
}

func (h *HTTPControllerCtx) SetRouter(r *Router) *HTTPControllerCtx {
	if r == nil {
		apperr.Throw(apperr.New("setRouter r cannot be nil!"))
	}

	h.router = r

	return h
}

// NewHTTPController return HTTPController
func NewHTTPController(arg ControllerArg) HTTPController {
	isDebug, _ := parser.String(os.Getenv("DEBUG")).ToBool()

	arg.Logger.Debug("OK", nil)

	h := &HTTPControllerCtx{
		Logger:         arg.Logger,
		ResponseMapper: arg.ResponseMapper,
		IsDebug:        isDebug,
		router:         arg.Router,
	}

	return h
}
