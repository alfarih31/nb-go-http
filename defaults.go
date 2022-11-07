package noob

import (
	"github.com/alfarih31/nb-go-http/utils"
	keyvalue "github.com/alfarih31/nb-go-keyvalue"
	"time"
)

var isDebug = false

var DefaultMeta = keyvalue.KeyValue{
	"app_name":        "Core",
	"app_description": "Core API",
	"app_version":     "v0.1.0",
}

var DefaultResponseHeader = map[string][]string{
	"Content-Type": {"application/json"},
}

var StartTime = utils.NewDatetimeNow().GetTime()

const defaultMaxBurstSize = 20
const defaultMaxEventPerSec = 1000

type CORSCfg struct {
	Enable           bool
	AllowOrigins     []string
	AllowMethods     string
	AllowHeaders     string
	AllowCredentials bool
	ExposeHeaders    string
	MaxAge           time.Duration
}

type ThrottlingCfg struct {
	Enable         bool
	MaxEventPerSec int
	MaxBurstSize   int
}

type Cfg struct {
	Host           string
	Port           int
	Path           string
	RequestTimeout time.Duration
	UseListener    bool
}

var DefaultCORSCfg = CORSCfg{
	Enable:           true,
	AllowOrigins:     nil, // use nil for wildcard origins (*)
	AllowHeaders:     "*",
	AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
	AllowCredentials: true,
	ExposeHeaders:    "authorization,content-type",
	MaxAge:           time.Duration(0),
}

var DefaultThrottlingCfg = ThrottlingCfg{
	Enable:         false,
	MaxBurstSize:   defaultMaxEventPerSec,
	MaxEventPerSec: defaultMaxBurstSize,
}

var DefaultCfg = Cfg{
	Host:           "",
	Port:           8080,
	Path:           "/",
	RequestTimeout: 0,
	UseListener:    false,
}

const (
	statusCodeOk uint = iota
	statusCodeErrInternal
	statusCodeErrNotFound
	statusCodeErrTooManyRequest
	statusCodeErrRequestTimeout
	statusCodeErrForbidden
)

var DefaultSuccessResponse = NewResponse(StatusOK, ResponseBody{
	Code:    statusCodeOk,
	Message: "success",
})

var DefaultSuccessNoContentResponse = NewResponse(StatusNoContent, ResponseBody{
	Code:    statusCodeOk,
	Message: "success",
})

var DefaultInternalServerErrorResponse = NewResponseError(StatusInternalServerError, ResponseBody{
	Code:    statusCodeErrInternal,
	Message: "internal server error",
})

var DefaultForbiddenErrorResponse = NewResponseError(StatusForbidden, ResponseBody{
	Code:    statusCodeErrForbidden,
	Message: "forbidden",
})

var DefaultNotFoundErrorResponse = NewResponseError(StatusNotFound, ResponseBody{
	Code:    statusCodeErrNotFound,
	Message: "not found",
})

var DefaultTooManyRequestsErrorResponse = NewResponseError(StatusTooManyRequests, ResponseBody{
	Code:    statusCodeErrTooManyRequest,
	Message: "too many request",
})

var DefaultRequestTimeoutErrorResponse = NewResponseError(StatusRequestTimeout, ResponseBody{
	Code:    statusCodeErrRequestTimeout,
	Message: "request timed out",
})
