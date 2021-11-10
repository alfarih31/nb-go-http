package main

import (
	"encoding/json"
	"github.com/alfarih31/nb-go-http/app_error"
	"github.com/alfarih31/nb-go-http/core"
	"github.com/alfarih31/nb-go-http/cors"
	"github.com/alfarih31/nb-go-http/data"
	"github.com/alfarih31/nb-go-http/env"
	http2 "github.com/alfarih31/nb-go-http/http"
	"github.com/alfarih31/nb-go-http/logger"
	"github.com/alfarih31/nb-go-http/utils"
	"net/http"
	"runtime/debug"
)

const (
	Success           = "OK"
	ErrorBadRequest   = "400"
	ErrorUnauthorized = "401"
	ErrorForbidden    = "403"
	ErrorNotFound     = "404"
	ErrorInternal     = "500"
	ErrorBadGateway   = "502"
)

var StandardResponses = map[string]http2.Response{
	Success: {
		Code: http.StatusOK,
		Body: ResponseBody{
			Status: ResponseStatus{
				Code:          0,
				MessageClient: "Success",
				MessageServer: "Success",
			},
		},
	},
	ErrorBadRequest: {
		Code: http.StatusBadRequest,
		Body: ResponseBody{
			Status: ResponseStatus{
				Code:          1,
				MessageClient: "Bad Request",
				MessageServer: "Bad Request",
			},
		},
	},
	ErrorUnauthorized: {
		Code: http.StatusUnauthorized,
		Body: ResponseBody{
			Status: ResponseStatus{
				Code:          1,
				MessageClient: "Unauthorized",
				MessageServer: "Unauthorized",
			},
		},
	},
	ErrorForbidden: {
		Code: http.StatusForbidden,
		Body: ResponseBody{
			Status: ResponseStatus{
				Code:          1,
				MessageClient: "Forbidden",
				MessageServer: "Forbidden",
			},
		},
	},
	ErrorNotFound: {
		Code: http.StatusNotFound,
		Body: ResponseBody{
			Status: ResponseStatus{
				Code:          1,
				MessageClient: "Not Found",
				MessageServer: "Not Found",
			},
		},
	},
	ErrorInternal: {
		Code: http.StatusInternalServerError,
		Body: ResponseBody{
			Status: ResponseStatus{
				Code:          1,
				MessageClient: "Internal Error",
				MessageServer: "Internal Error",
			},
		},
	},
	ErrorBadGateway: {
		Code: http.StatusBadGateway,
		Body: ResponseBody{
			Status: ResponseStatus{
				Code:          1,
				MessageClient: "Bad Gateway",
				MessageServer: "Bad Gateway",
			},
		},
	},
}

type ResponseHeader = data.KeyValue

type ResponseStatus struct {
	Code          int    `json:"code"`
	MessageClient string `json:"message_client"`
	MessageServer string `json:"message_server"`
}

func (s ResponseStatus) AssignTo(target *ResponseStatus) {
	if target.MessageClient == "" {
		target.MessageClient = s.MessageClient
	}

	if target.MessageServer == "" {
		target.MessageServer = s.MessageServer
	}

	if target.Code == 0 {
		target.Code = s.Code
	}
}

type ResponseBody struct {
	Status ResponseStatus `json:"status"`
	Meta   interface{}    `json:"meta,omitempty"`
	Data   interface{}    `json:"data,omitempty"`
	Errors []interface{}  `json:"errors,omitempty"`
}

func (b ResponseBody) Compose(body ResponseBody) ResponseBody {

	b.Status.AssignTo(&body.Status)

	if body.Errors == nil {
		body.Errors = b.Errors
	}

	return body
}

func (b ResponseBody) String() string {
	j, err := json.Marshal(b)

	if err != nil {
		return ""
	}

	return string(j)
}

func main() {
	env, _ := env.LoadEnv(".env")

	rl := logger.Logger("RootLogger")

	basePath, _ := env.GetString("SERVER_PATH", "/v1")
	baseHost, _ := env.GetString("SERVER_HOST", ":")
	basePort, _ := env.GetInt("SERVER_PORT", 8080)

	utils.Func(utils.Run{
		Try: func() {
			responseMapper := http2.ResponseMapper(http2.ResponseMapperCfg{
				Logger:            logger.Logger("ResponseMapper"),
				SuccessCode:       "OK",
				InternalErrorCode: "500",
			})
			responseMapper.Load(StandardResponses)

			app := core.Core(&core.CoreCfg{
				Meta: &data.KeyValue{
					"app_name":        "test",
					"app_version":     "v0.1.0",
					"app_description": "Description",
				},
			})

			app.Setup = func() {
				g1 := http2.HTTPController(app.RootController.BranchRouter("/sample"), app.Logger.NewChild("G1-Controller"), responseMapper)
				g2 := http2.HTTPController(g1.BranchRouter("/deep"), app.Logger.NewChild("G2-Controller"), responseMapper)

				g1.Handle("GET /first-inner", func(c *http2.HandlerCtx) *http2.Response {
					return &http2.Response{
						Body: ResponseBody{
							Status: ResponseStatus{
								MessageClient: "G1 FIRST",
							},
						},
					}
				})

				g1.Handle("GET /error", func(c *http2.HandlerCtx) *http2.Response {
					http2.HTTPError.BadGateway.Throw(nil)
					return nil
				})

				g1.Handle("GET /second-inner", func(c *http2.HandlerCtx) *http2.Response {
					return &http2.Response{
						Body: ResponseBody{
							Status: ResponseStatus{
								MessageClient: "G1 SECOND",
							},
							Data: []string{
								"1", "2", "3",
							},
						},
					}
				})

				g2.Handle("GET /first-inner", func(context *http2.HandlerCtx) *http2.Response {
					return &http2.Response{
						Body: ResponseBody{
							Status: ResponseStatus{
								MessageClient: "G2 FIRST",
							},
							Data: []string{
								"1", "2", "3",
							},
						},
					}
				})

				app.Logger.Debug("Init Controllers OK...", nil)
			}

			app.Start(core.StartArg{
				Host:           baseHost,
				Path:           basePath,
				Port:           basePort,
				ResponseMapper: &responseMapper,
				CORS: &cors.CORSCfg{
					Enable: true,
				},
			})
		},
		Catch: func(e interface{}) {
			ee, ok := e.(apperror.Err)

			debug.PrintStack()
			if ok {
				rl.Error(ee, map[string]interface{}{"error": ee.Errors(), "stack": apperror.StackTrace()})
			} else {
				rl.Error(ee, nil)
			}
		},
	})
}
