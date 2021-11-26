module github.com/alfarih31/nb-go-http/tcf

replace (
	github.com/alfarih31/nb-go-http/app_err => ../app_err
	github.com/alfarih31/nb-go-http/keyvalue => ../keyvalue
)

go 1.17

require github.com/alfarih31/nb-go-http/app_err v0.0.0-00010101000000-000000000000

require (
	github.com/DataDog/gostackparse v0.5.0 // indirect
	github.com/alfarih31/nb-go-http/keyvalue v0.0.0-00010101000000-000000000000 // indirect
)
