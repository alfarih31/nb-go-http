// This file is not part of the project structure. This file is just an example

package main

import (
	"encoding/json"
	"fmt"
	"github.com/alfarih31/nb-go-http"
	"github.com/alfarih31/nb-go-http/app_err"
	"github.com/alfarih31/nb-go-http/cors"
	_env "github.com/alfarih31/nb-go-http/env"
	"github.com/alfarih31/nb-go-http/logger"
	"github.com/alfarih31/nb-go-http/tcf"
	"github.com/alfarih31/nb-go-keyvalue"
	"net/http"
	"runtime"
	"runtime/debug"
)

type Standard struct {
	Success            noob.Response
	ErrorInternal      noob.Response
	ErrorBadRequest    noob.Response
	ErrorNotFound      noob.Response
	ErrorToManyRequest noob.Response
}

var StandardResponses = Standard{
	Success: &noob.DefaultResponse{
		Code: http.StatusOK,
		Body: ResponseBody{
			Status: ResponseStatus{
				Code:          0,
				MessageClient: "Success",
				MessageServer: "Success",
			},
		},
	},
	ErrorBadRequest: &noob.DefaultResponse{
		Code: http.StatusBadRequest,
		Body: ResponseBody{
			Status: ResponseStatus{
				Code:          1,
				MessageClient: "Bad Request",
				MessageServer: "Bad Request",
			},
		},
	},
	ErrorNotFound: &noob.DefaultResponse{
		Code: http.StatusNotFound,
		Body: ResponseBody{
			Status: ResponseStatus{
				Code:          1,
				MessageClient: "Not Found",
				MessageServer: "Not Found",
			},
		},
	},
	ErrorInternal: &noob.DefaultResponse{
		Code: http.StatusInternalServerError,
		Body: ResponseBody{
			Status: ResponseStatus{
				Code:          1,
				MessageClient: "Internal Error",
				MessageServer: "Internal Error",
			},
		},
	},
	ErrorToManyRequest: &noob.DefaultResponse{
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
			})

			app := noob.New(&noob.CoreCfg{
				ResponseMapper: responseMapper,
				Meta: keyvalue.KeyValue{
					"app_name":        "test",
					"app_version":     "v0.1.0",
					"app_description": "Description",
				},
			})

			app.Setup = func() error {

				g1 := noob.NewHTTPController(noob.ControllerArg{
					Logger:         app.Logger.NewChild("G1-Controller"),
					ResponseMapper: responseMapper,
				}).SetRouter(app.BranchRouter("/sample"))

				g2 := noob.NewHTTPController(noob.ControllerArg{
					Logger:         app.Logger.NewChild("G2-Controller"),
					ResponseMapper: responseMapper,
				}).SetRouter(g1.BranchRouter("/deep"))

				g1.Handle("GET /first-inner", func(c *noob.HandlerCtx) (noob.Response, error) {
					qp := noob.QueryParser(*c)
					q1, err := qp.GetInt("q1")
					fmt.Println("q1", q1, err)

					qReq, err := qp.GetInt("qreq", noob.QueryParserOption{
						Required: true,
					})
					fmt.Println("qreq", qReq, err)

					qReqWithDef, err := qp.GetInt("qreqdef", noob.QueryParserOption{
						Default:  345,
						Required: true,
					})
					fmt.Println("qreqWithDef", qReqWithDef, err)

					return &noob.DefaultResponse{
						Body: ResponseBody{
							Status: ResponseStatus{
								MessageClient: "G1 FIRST",
							},
						},
					}, nil
				})

				g1.Handle("GET /error", func(c *noob.HandlerCtx) (noob.Response, error) {
					panic("THIS IS AN ERROR")
				})

				g1.Handle("GET /second-inner", func(c *noob.HandlerCtx) (noob.Response, error) {
					res := StandardResponses.Success
					res.SetBody(ResponseBody{
						Status: ResponseStatus{
							MessageClient: "G1 SECOND",
						},
						Data: "G1 SECOND DATA",
					})

					return res, nil
				})

				g2.Handle("GET /first-inner", func(context *noob.HandlerCtx) (noob.Response, error) {
					return &noob.DefaultResponse{
						Body: ResponseBody{
							Status: ResponseStatus{
								MessageClient: "G2 FIRST",
							},
							Data: []string{
								"1", "2", "3",
							},
						},
					}, nil
				})

				app.Logger.Debug("Init Controllers OK...", nil)

				return nil
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
