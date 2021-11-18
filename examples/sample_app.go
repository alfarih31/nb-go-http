package main

import (
	"encoding/json"
	"github.com/alfarih31/nb-go-http"
	"github.com/alfarih31/nb-go-http/app_err"
	"github.com/alfarih31/nb-go-http/cors"
	_env "github.com/alfarih31/nb-go-http/env"
	"github.com/alfarih31/nb-go-http/keyvalue"
	"github.com/alfarih31/nb-go-http/logger"
	"github.com/alfarih31/nb-go-http/tcf"
	"net/http"
	"runtime"
	"runtime/debug"
)

const (
	Success            = "OK"
	ErrorBadRequest    = "400"
	ErrorNotFound      = "404"
	ErrorInternal      = "500"
	ErrorToManyRequest = "429"
)

type Response = noob.Response

type Standard struct {
	Success            Response
	ErrorInternal      Response
	ErrorBadRequest    Response
	ErrorNotFound      Response
	ErrorToManyRequest Response
}

var StandardResponses = Standard{
	Success: Response{
		Code: http.StatusOK,
		Body: ResponseBody{
			Status: ResponseStatus{
				Code:          0,
				MessageClient: "Success",
				MessageServer: "Success",
			},
		},
	},
	ErrorBadRequest: Response{
		Code: http.StatusBadRequest,
		Body: ResponseBody{
			Status: ResponseStatus{
				Code:          1,
				MessageClient: "Bad Request",
				MessageServer: "Bad Request",
			},
		},
	},
	ErrorNotFound: Response{
		Code: http.StatusNotFound,
		Body: ResponseBody{
			Status: ResponseStatus{
				Code:          1,
				MessageClient: "Not Found",
				MessageServer: "Not Found",
			},
		},
	},
	ErrorInternal: Response{
		Code: http.StatusInternalServerError,
		Body: ResponseBody{
			Status: ResponseStatus{
				Code:          1,
				MessageClient: "Internal Error",
				MessageServer: "Internal Error",
			},
		},
	},
	ErrorToManyRequest: Response{
		Code: http.StatusTooManyRequests,
		Body: ResponseBody{
			Status: ResponseStatus{
				Code:          1,
				MessageClient: "Too Many Request",
				MessageServer: "Too Many Request",
			},
		},
	},
}

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
	env, _ := _env.LoadEnv(".env")

	rl := logger.New("RootLogger")

	basePath, _ := env.GetString("SERVER_PATH", "/v1")
	baseHost, _ := env.GetString("SERVER_HOST", "")
	basePort, _ := env.GetInt("SERVER_PORT", 8080)

	tcf.TCFunc(tcf.Func{
		Try: func() {
			responseMapper := noob.ResponseMapper(noob.ResponseMapperCfg{
				Logger: logger.New("ResponseMapper"),
				DefaultCode: &noob.DefaultResponseCode{
					InternalError: ErrorInternal,
				},
			})

			responseMapper.Load(map[string]noob.Response{
				Success:            StandardResponses.Success,
				ErrorInternal:      StandardResponses.ErrorInternal,
				ErrorBadRequest:    StandardResponses.ErrorBadRequest,
				ErrorToManyRequest: StandardResponses.ErrorToManyRequest,
			})

			app := noob.New(&noob.CoreCfg{
				ResponseMapper: responseMapper,
				Meta: &keyvalue.KeyValue{
					"app_name":        "test",
					"app_version":     "v0.1.0",
					"app_description": "Description",
				},
			})

			app.Setup = func() {

				app.Handle("USE", func(context *noob.HandlerCtx) *noob.Response {
					context.Next()
					return nil
				})

				g1 := noob.NewController(noob.ControllerArg{
					Logger:         app.Logger.NewChild("G1-Controller"),
					ResponseMapper: responseMapper,
				}).SetRouter(app.BranchRouter("/sample"))

				g2 := noob.NewController(noob.ControllerArg{
					Logger:         app.Logger.NewChild("G2-Controller"),
					ResponseMapper: responseMapper,
				}).SetRouter(g1.BranchRouter("/deep"))

				g1.Handle("GET /first-inner", func(c *noob.HandlerCtx) *noob.Response {
					return &noob.Response{
						Body: ResponseBody{
							Status: ResponseStatus{
								MessageClient: "G1 FIRST",
							},
						},
					}
				})

				g1.Handle("GET /error", func(c *noob.HandlerCtx) *noob.Response {
					panic("THIS IS AN ERROR")
					return nil
				})

				g1.Handle("GET /second-inner", func(c *noob.HandlerCtx) *noob.Response {
					res := StandardResponses.Success
					res.Body = ResponseBody{
						Status: ResponseStatus{
							MessageClient: "G1 SECOND",
						},
						Data: "G1 SECOND DATA",
					}

					return &res
				})

				g2.Handle("GET /first-inner", func(context *noob.HandlerCtx) *noob.Response {
					return &noob.Response{
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

			app.Start(noob.StartArg{
				Host: baseHost,
				Path: basePath,
				Port: basePort,
				CORS: &cors.Cfg{
					Enable: true,
				},
				Throttling: &noob.ThrottlingCfg{
					MaxEventPerSec: 2,
					MaxBurstSize:   1,
				},
			})
		},
		Catch: func(e interface{}, frames *runtime.Frames) {
			ee, ok := e.(apperr.AppErr)

			debug.PrintStack()
			if ok {
				rl.Error(ee, map[string]interface{}{"error": ee.Errors(), "stack": apperr.StackTrace()})
			} else {
				rl.Error(ee, nil)
			}
		},
	})
}
