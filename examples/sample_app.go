package main

import (
	"encoding/json"
	"github.com/alfarih31/nb-go-http"
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

var StandardResponses = map[string]nb_go.Response{
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
	env, _ := nb_go.LoadEnv(".env")

	rl := nb_go.Logger("RootLogger")

	basePath, _ := env.GetString("SERVER_PATH", "/v1")
	baseHost, _ := env.GetString("SERVER_HOST", ":")
	basePort, _ := env.GetInt("SERVER_PORT", 8080)

	nb_go.FlowFunc(nb_go.Func{
		Try: func() {
			responseMapper := nb_go.ResponseMapper(nb_go.ResponseMapperCfg{
				Logger:            nb_go.Logger("ResponseMapper"),
				SuccessCode:       "OK",
				InternalErrorCode: "500",
			})
			responseMapper.Load(StandardResponses)

			app := nb_go.Core(&nb_go.CoreCfg{
				Meta: &nb_go.KeyValue{
					"app_name":        "test",
					"app_version":     "v0.1.0",
					"app_description": "Description",
				},
			})

			app.Setup = func() {
				g1 := nb_go.HTTPController(app.RootController.BranchRouter("/sample"), app.Logger.NewChild("G1-Controller"), responseMapper)
				g2 := nb_go.HTTPController(g1.BranchRouter("/deep"), app.Logger.NewChild("G2-Controller"), responseMapper)

				g1.Handle("GET /first-inner", func(c *nb_go.HandlerCtx) *nb_go.Response {
					return &nb_go.Response{
						Body: ResponseBody{
							Status: ResponseStatus{
								MessageClient: "G1 FIRST",
							},
						},
					}
				})

				g1.Handle("GET /error", func(c *nb_go.HandlerCtx) *nb_go.Response {
					nb_go.HTTPError.BadGateway.Throw(nil)
					return nil
				})

				g1.Handle("GET /second-inner", func(c *nb_go.HandlerCtx) *nb_go.Response {
					return &nb_go.Response{
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

				g2.Handle("GET /first-inner", func(context *nb_go.HandlerCtx) *nb_go.Response {
					return &nb_go.Response{
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

			app.Start(nb_go.StartArg{
				Host:           baseHost,
				Path:           basePath,
				Port:           basePort,
				ResponseMapper: &responseMapper,
				CORS: &nb_go.CORSCfg{
					Enable: true,
				},
			})
		},
		Catch: func(e interface{}) {
			ee, ok := e.(nb_go.Err)

			debug.PrintStack()
			if ok {
				rl.Error(ee, map[string]interface{}{"error": ee.Errors(), "stack": nb_go.StackTrace()})
			} else {
				rl.Error(ee, nil)
			}
		},
	})
}
