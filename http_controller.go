package nbgohttp

import (
	"fmt"
    "math"
    "net/http"
	"strings"
)

const abortIndex int8 = math.MaxInt8 / 2

type HTTPControllerCtx struct {
	Router         *ExtRouter
	Logger         ILogger
	ResponseMapper IResponseMapper
	Debug          bool
}

type HandlerSpec struct {
	method string
	path   string
}

func (h HTTPControllerCtx) GetSpec(spec string) HandlerSpec {
	if spec == "" {
		ThrowError(Err{Message: "Handler Spec cannot be empty!"})
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
		Func(FuncCtx{
			Try: func() {
				res := handler(c)

				if res != nil {
					if res.Body.Errors != nil {
						h.SendError(c, NewError(Err{
							Message: "Request Catch an Error",
							Err:     res.Body.Errors,
						}))

						return
					}

					h.SendSuccess(c, res)
				}
			},
			Catch: func(e interface{}) {
				h.Logger.Debug(e, nil)
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

func (h *HTTPControllerCtx) ChainControllers(ctrls ...HTTPControllerCtx) {
    for _, ctrl := range ctrls {
        h.Router.AppendHandlers(ctrl.Router.Handlers())
    }
}

func (h *HTTPControllerCtx) Handle(spec string, handlers ...HTTPHandler) {
	if &h.Router == nil {
		ThrowError(Err{Message: "Cannot Set Spec, router is nil"})
	}

	handlerSpec := h.GetSpec(spec)

	switch strings.ToUpper(handlerSpec.method) {
	case "GET":
		h.Router.GET(handlerSpec.path, h.ToExtHandlers(handlers)...)
	case "POST":
		h.Router.POST(handlerSpec.path, h.ToExtHandlers(handlers)...)
	case "PUT":
		h.Router.PUT(handlerSpec.path, h.ToExtHandlers(handlers)...)
	case "DELETE":
		h.Router.DELETE(handlerSpec.path, h.ToExtHandlers(handlers)...)
	case "OPTIONS":
		h.Router.OPTIONS(handlerSpec.path, h.ToExtHandlers(handlers)...)
	case "HEAD":
		h.Router.HEAD(handlerSpec.path, h.ToExtHandlers(handlers)...)
    case "PATCH":
        h.Router.PATCH(handlerSpec.path, h.ToExtHandlers(handlers)...)
	case "USE":
		h.Router.USE(h.ToExtHandlers(handlers)...)
	default:
		ThrowError(Err{Message: fmt.Sprintf("Unknown HTTP Handle Spec: %s", spec)})
	}
}

func (h *HTTPControllerCtx) SendSuccess(c *HandlerCtx, res *Response) {
	r := h.ResponseMapper.GetSuccess()

    r = r.Compose(*res)
    
	_, e := c.Response(http.StatusOK, r.Body.String(), r.Header)

	if e != nil {
		h.Logger.Debug("", map[string]interface{}{"_error": e})
	}
}

func (h *HTTPControllerCtx) SendError(c *HandlerCtx, e interface{}) {
	errorData := KeyValue{}

	r := h.ResponseMapper.GetInternalError()

	switch er := e.(type) {
	case Err:
		r = h.ResponseMapper.Get(er.Code, nil)

        errorData["errors"] = er.Data

		if er.Code != "" && er.Code != StatusErrorInternal && r.Code == http.StatusInternalServerError {
			h.Logger.Warn(fmt.Sprintf("Error code not mapped. Code = %s", er.Code), nil)
		}
	case error:
		r.Body = r.Body.Compose(ResponseBody{
        Status: ResponseStatus{
            MessageServer: er.Error(),
        },
        })
	case string:
        r.Body = r.Body.Compose(ResponseBody{
            Status: ResponseStatus{
                MessageServer: er,
            },
        })
	}

	_, rEr := c.Response(r.Code, r.Body.String(), r.Header)

	if rEr != nil {
		h.Logger.Debug("", map[string]interface{}{"_error": rEr})
	}
}

func (h *HTTPControllerCtx) SetRouter(r ExtRouter) {
	if &r == nil {
		ThrowError(Err{Message: "setRouter r cannot be nil!"})
	}

	h.Router = &r
}

func HTTPController(r *ExtRouter, l ILogger, rm IResponseMapper) HTTPControllerCtx {
	l.Debug("OK", nil)

	h := HTTPControllerCtx{
		Router:         r,
		Logger:         l,
		ResponseMapper: rm,
	}
    
	return h
}
