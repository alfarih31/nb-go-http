package nbgohttp

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"os"
	"strings"
)

const abortIndex int8 = math.MaxInt8 / 2

type HTTPControllerCtx struct {
	Context        context.Context
	router         *ExtRouter
	Logger         ILogger
	ResponseMapper *ResponseMapperCtx
	Debug          bool
}

type HTTPControllerArg struct {
	Context        context.Context
	Router         *ExtRouter
	Logger         ILogger
	ResponseMapper *ResponseMapperCtx
}

type HandlerSpec struct {
	method string
	path   string
}

func (h HTTPControllerCtx) GetSpec(spec string) HandlerSpec {
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

func (h HTTPControllerCtx) ToExtHandler(handler HTTPHandler) ExtHandler {
	return func(ec *ExtHandlerCtx) {
		c := WrapExtHandlerCtx(ec)
		FlowFunc(Func{
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

func (h HTTPControllerCtx) ToExtHandlers(handlers []HTTPHandler) []ExtHandler {
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
	if &h.router == nil {
		ThrowError(&Err{Message: "Cannot Set Spec, router is nil"})
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
	case *Err:
		r = h.ResponseMapper.Get(er.Code, nil)

		if !h.Debug {
			er.Stack = nil
		}

		if r.Code == http.StatusInternalServerError {
			r.ComposeBody(KeyValue{"errors": er})
		}
	case Err:
		r = h.ResponseMapper.Get(er.Code, nil)

		if !h.Debug {
			er.Stack = nil
		}

		if r.Code == http.StatusInternalServerError {
			r.ComposeBody(KeyValue{"errors": er})
		}
	case Response:
		r = er
	case *Response:
		r = *er
	case error:
		r.ComposeBody(KeyValue{"errors": er})
	case string:
		r.ComposeBody(KeyValue{"errors": er})
	}

	_, rEr := c.response(r.Code, r.Body, r.Header)

	if rEr != nil {
		h.Logger.Debug("", map[string]interface{}{"_error": rEr})
	}
}

func (h *HTTPControllerCtx) SetRouter(r *ExtRouter) {
	if r == nil {
		ThrowError(&Err{Message: "setRouter r cannot be nil!"})
	}

	h.router = r
}

func (h *HTTPControllerCtx) WithContext(ctx context.Context) *HTTPControllerCtx {
	return HTTPController(HTTPControllerArg{
		Context:        ctx,
		Router:         h.router,
		Logger:         h.Logger,
		ResponseMapper: h.ResponseMapper,
	})
}

func HTTPController(arg HTTPControllerArg) *HTTPControllerCtx {
	isDebug, _ := StringParser{os.Getenv("DEBUG")}.ToBool()

	arg.Logger.Debug("OK", nil)

	h := &HTTPControllerCtx{
		Context:        arg.Context,
		router:         arg.Router,
		Logger:         arg.Logger,
		ResponseMapper: arg.ResponseMapper,
		Debug:          isDebug,
	}

	if arg.Context != nil && arg.ResponseMapper != nil {
		arg.ResponseMapper = arg.ResponseMapper.WithContext(arg.Context)
	}

	return h
}
