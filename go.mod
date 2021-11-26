module github.com/alfarih31/nb-go-http

replace (
	github.com/alfarih31/nb-go-http/app_err => ./app_err
	github.com/alfarih31/nb-go-http/env => ./env
	github.com/alfarih31/nb-go-http/keyvalue => ./keyvalue
	github.com/alfarih31/nb-go-http/logger => ./logger
	github.com/alfarih31/nb-go-http/parser => ./parser
	github.com/alfarih31/nb-go-http/tcf => ./tcf
)

go 1.17

require (
	github.com/alfarih31/nb-go-http/app_err v1.3.19
	github.com/alfarih31/nb-go-http/env v1.3.19
	github.com/alfarih31/nb-go-http/keyvalue v1.3.19
	github.com/alfarih31/nb-go-http/logger v1.3.19
	github.com/alfarih31/nb-go-http/parser v1.3.19
	github.com/alfarih31/nb-go-http/tcf v1.3.19
	github.com/gin-gonic/gin v1.7.4
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac
)

require (
	github.com/DataDog/gostackparse v0.5.0 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator/v10 v10.9.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/ugorji/go/codec v1.2.6 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/sys v0.0.0-20211103235746-7861aae1554b // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
