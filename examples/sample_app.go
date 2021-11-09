package main

import (
	"github.com/alfarih31/nb-go-http"
	"runtime/debug"
)

func main() {
	env, _ := nbgohttp.LoadEnv(".env")

	rl := nbgohttp.Logger("RootLogger")

	isDebug, _ := env.GetBool("DEBUG", false)
	basePath, _ := env.GetString("SERVER_PATH", "/")
	baseHost, _ := env.GetString("SERVER_HOST", ":")
	basePort, _ := env.GetInt("SERVER_PORT", 8080)

	httpErrs := nbgohttp.HttpError()

	nbgohttp.Func(nbgohttp.FuncCtx{
		Try: func() {
			app := nbgohttp.Core(&nbgohttp.CoreCfg{
				Debug: isDebug,
				Server: &nbgohttp.ServerCfg{
					Host: baseHost,
					Path: basePath,
					Port: basePort,
					CORS: &nbgohttp.CORSCfg{
						Enable: false,
					},
				},
				Meta: &nbgohttp.KeyValue{
					"app_name":        "test",
					"app_version":     "v0.1.0",
					"app_description": "Description",
				},
			})

			app.Setup = func() {
				g := nbgohttp.HTTPController(app.Router.Branch("/test"), app.Logger.NewChild("TestController"), app.HTTP.ResponseMapper)

				g.Handle("GET /first-inner", func(c *nbgohttp.HandlerCtx) *nbgohttp.Response {
					return &nbgohttp.Response{
						Body: nbgohttp.ResponseBody{
							Status: nbgohttp.ResponseStatus{
								MessageClient: "FIRST",
							},
						},
					}
				})

				g.Handle("GET /error", func(c *nbgohttp.HandlerCtx) *nbgohttp.Response {
					httpErrs.BadGateway.Throw(nil)
					return nil
				})

				g.Handle("GET /second-inner", func(c *nbgohttp.HandlerCtx) *nbgohttp.Response {
					return &nbgohttp.Response{
						Body: nbgohttp.ResponseBody{
							Status: nbgohttp.ResponseStatus{
								MessageClient: "SECOND",
							},
							Data: []string{
								"1", "2", "3",
							},
						},
					}
				})

				app.HTTP.ChainControllers(g)

				app.Logger.Debug("Init Controllers OK...", nil)
			}

			app.Start()
		},
		Catch: func(e interface{}) {
			ee, ok := e.(nbgohttp.Err)

			debug.PrintStack()
			if ok {
				rl.Error(ee, map[string]interface{}{"error": ee.Errors(), "stack": ee.Stack()})
			} else {
				rl.Error(ee, nil)
			}
		},
	})
}
