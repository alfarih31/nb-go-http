// This file is not part of the project structure. This file is just an example

package main

import (
	"encoding/json"
	"errors"
	_env "github.com/alfarih31/nb-go-env"
	"github.com/alfarih31/nb-go-http"
	"github.com/alfarih31/nb-go-keyvalue"
	logger "github.com/alfarih31/nb-go-logger"
	"runtime"
)

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

	noob.TCFunc(noob.Func{
		Try: func() {
			noob.DefaultMeta = keyvalue.KeyValue{
				"app_name":        "test",
				"app_version":     "v0.1.0",
				"app_description": "Description",
			}

			noob.DefaultCfg = noob.Cfg{
				Host: baseHost,
				Path: basePath,
				Port: basePort,
			}

			noob.DefaultThrottlingCfg = noob.ThrottlingCfg{
				MaxEventPerSec: 2,
				MaxBurstSize:   1,
			}

			noob.DefaultCORSCfg = noob.CORSCfg{
				Enable:       true,
				AllowOrigins: []string{"*"},
			}

			app := noob.New()

			g1 := app.Branch("/sample")

			g2 := g1.Branch("/deep")

			g1.GET("/first-inner", func(c *noob.HandlerCtx) (noob.Response, error) {

				return noob.NewResponseSuccess(noob.ResponseBody{
					Message: "G1 First",
				}), nil
			})

			g1.GET("/error", func(c *noob.HandlerCtx) (noob.Response, error) {
				return nil, errors.New("this is an error")
			})

			g1.GET("/second-inner", func(c *noob.HandlerCtx) (noob.Response, error) {
				return noob.NewResponseSuccess(noob.ResponseBody{
					Message: "G1 Second",
				}), nil
			})

			g2.GET("/first-inner", func(context *noob.HandlerCtx) (noob.Response, error) {
				return noob.NewResponseSuccess(noob.ResponseBody{
					Message: "G2 First",
					Data: []string{
						"1", "2", "3",
					},
				}), nil
			})

			if err := app.Start(); err != nil {
				panic(err)
			}
		},
		Catch: func(e interface{}, frames *runtime.Frames) {
			rl.Error(e)
		},
	})
}
