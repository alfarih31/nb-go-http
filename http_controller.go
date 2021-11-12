package noob

import (
	"fmt"
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

type httpControllerCtx struct {
	router         *ExtRouter
	Logger         logger.ILogger
	ResponseMapper *responseMapperCtx
	Debug          bool
}

type ControllerArg struct {
	Logger         logger.ILogger
	ResponseMapper *responseMapperCtx
}

type HandlerSpec struct {
	method string
	path   string
}

func (h *httpControllerCtx) GetSpec(spec string) HandlerSpec {
	if spec == "" {
		ThrowError(&Err{Message: "Handler Spec cannot be empty!"})
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

func (h *httpControllerCtx) ToExtHandler(handler HTTPHandler) ExtHandler {
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

func (h *httpControllerCtx) ToExtHandlers(handlers []HTTPHandler) []ExtHandler {
	extHandlers := make([]ExtHandler, len(handlers), len(handlers))

	for i, handler := range handlers {
		extHandlers[i] = h.ToExtHandler(handler)
	}

	return extHandlers
}

func (h *httpControllerCtx) BranchRouter(path string) *ExtRouter {
	h.Logger.Debug(fmt.Sprintf("Branching Controller Router with Path : %s", path), map[string]interface{}{
		"branchFullPath": fmt.Sprintf("%s%s", h.router.fullPath, path),
	})

	return h.router.Branch(path, fmt.Sprintf("%s%s", h.router.fullPath, path))
}

func (h *httpControllerCtx) Handle(spec string, handlers ...HTTPHandler) {
	if h.router == nil {
		ThrowError(&Err{Message: "Cannot Set Spec, SetRouter first!"})
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
		ThrowError(&Err{Message: fmt.Sprintf("Unknown HTTP Handle Spec: %s", spec)})
	}
}

func (h *httpControllerCtx) SendSuccess(c *HandlerCtx, res *Response) {
	r := h.ResponseMapper.GetSuccess()

	r.ComposeTo(res)

	_, e := c.response(res.Code, res.Body, res.Header)

	if e != nil {
		h.Logger.Debug("", map[string]interface{}{"_error": e})
	}
}

func (h *httpControllerCtx) SendError(c *HandlerCtx, e interface{}) {
	r := h.ResponseMapper.GetInternalError()

	switch er := e.(type) {
	case *Err:
		r = h.ResponseMapper.Get(er.Code, nil)

		if !h.Debug {
			er.Stack = nil
		}

		if r.Code == http.StatusInternalServerError {
			r.ComposeBody(keyvalue.KeyValue{"errors": er})
		}
	case Err:
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
	}

	_, rEr := c.response(r.Code, r.Body, r.Header)

	if rEr != nil {
		h.Logger.Debug("", map[string]interface{}{"_error": rEr})
	}
}

func (h *httpControllerCtx) SetRouter(r *ExtRouter) *httpControllerCtx {
	if r == nil {
		ThrowError(&Err{Message: "setRouter r cannot be nil!"})
	}

	h.router = r

	return h
}

func NewController(arg ControllerArg) *httpControllerCtx {
	isDebug, _ := parser.String(os.Getenv("DEBUG")).ToBool()

	arg.Logger.Debug("OK", nil)

	h := &httpControllerCtx{
		Logger:         arg.Logger,
		ResponseMapper: arg.ResponseMapper,
		Debug:          isDebug,
	}

	return h
}
