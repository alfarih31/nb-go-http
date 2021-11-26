module github.com/alfarih31/nb-go-http/tcf

replace (
	github.com/alfarih31/nb-go-http/app_err => ../app_err
	github.com/alfarih31/nb-go-http/keyvalue => ../keyvalue
)

go 1.17

require github.com/alfarih31/nb-go-http/app_err v1.3.19

require (
	github.com/DataDog/gostackparse v0.5.0 // indirect
	github.com/alfarih31/nb-go-http/keyvalue v1.3.19 // indirect
)
