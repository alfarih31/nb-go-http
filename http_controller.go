package noob

import (
	"fmt"
	"github.com/alfarih31/nb-go-http/app_err"
	"github.com/alfarih31/nb-go-http/keyvalue"
	"github.com/alfarih31/nb-go-http/logger"
	"github.com/alfarih31/nb-go-http/parser"
	"github.com/alfarih31/nb-go-http/tcf"
	"math"
	"net/http"
	"os"
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
			Catch: func(e interface{}) {
				h.SendError(c, e)
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

	_, e := c.response(res.Code, res.Body, res.Header)

	if e != nil {
		h.Logger.Debug("", map[string]interface{}{"_error": e})
	}
}

func (h *HTTPControllerCtx) SendError(c *HandlerCtx, e interface{}) {
	r := h.ResponseMapper.GetInternalError()

	switch er := e.(type) {
	case *apperr.AppErr:
		r = h.ResponseMapper.Get(er.Code, nil)

		if !h.Debug {
			er.Stack = nil
		}

		if r.Code == http.StatusInternalServerError {
			r.ComposeBody(keyvalue.KeyValue{"errors": er})
		}
	case apperr.AppErr:
		r = h.ResponseMapper.Get(er.Code, nil)

		if !h.Debug {
			er.Stack = nil
		}

		if r.Code == http.StatusInternalServerError {
			r.ComposeBody(keyvalue.KeyValue{"errors": er})
		}
	case Response:
		r = er
	case *Response:
		r = *er
	case error:
		r.ComposeBody(keyvalue.KeyValue{"errors": er})
	case string:
		r.ComposeBody(keyvalue.KeyValue{"errors": er})
	default:
		r.ComposeBody(keyvalue.KeyValue{"errors": er})
	}

	_, rEr := c.response(r.Code, r.Body, r.Header)

	if rEr != nil {
		h.Logger.Debug("", map[string]interface{}{"_error": rEr})
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
