package noob

import (
	"fmt"
	"github.com/alfarih31/nb-go-http/app_err"
	"github.com/alfarih31/nb-go-http/logger"
	"github.com/alfarih31/nb-go-http/parser"
	"github.com/alfarih31/nb-go-keyvalue"
	"net/http"
	"reflect"
)

type ResponseMapperCfg struct {
	Logger      logger.Logger
	DefaultCode *DefaultResponseCode
}

type DefaultResponseCode struct {
	Success       string
	InternalError string
}

type _builtInResponse struct {
	ResponseSuccess        Response
	ResponseBadRequest     Response
	ResponseUnauthorized   Response
	ResponseForbidden      Response
	ResponseInternalError  Response
	ResponseBadGateway     Response
	ResponseNotFound       Response
	ResponseNoMethod       Response
	ResponseTooManyRequest Response
}

var builtInResponse = _builtInResponse{
	ResponseSuccess: Response{
		Code: http.StatusOK,
		Body: keyvalue.KeyValue{},
		Header: map[string]string{
			"Content-Type": "application/json",
		},
	},
	ResponseBadRequest: Response{
		Code: http.StatusBadRequest,
		Body: keyvalue.KeyValue{},
	},
	ResponseUnauthorized: Response{
		Code: http.StatusUnauthorized,
		Body: keyvalue.KeyValue{},
		Header: map[string]string{
			"Content-Type": "application/json",
		},
	},
	ResponseForbidden: Response{
		Code: http.StatusForbidden,
		Body: keyvalue.KeyValue{},
		Header: map[string]string{
			"Content-Type": "application/json",
		},
	},
	ResponseInternalError: Response{
		Code: http.StatusInternalServerError,
		Body: keyvalue.KeyValue{},
		Header: map[string]string{
			"Content-Type": "application/json",
		},
	},
	ResponseBadGateway: Response{
		Code: http.StatusBadGateway,
		Body: keyvalue.KeyValue{},
		Header: map[string]string{
			"Content-Type": "application/json",
		},
	},
	ResponseNotFound: Response{
		Code: http.StatusNotFound,
		Body: keyvalue.KeyValue{},
		Header: map[string]string{
			"Content-Type": "application/json",
		},
	},
	ResponseNoMethod: Response{
		Code: http.StatusMethodNotAllowed,
		Body: keyvalue.KeyValue{},
		Header: map[string]string{
			"Content-Type": "application/json",
		},
	},
	ResponseTooManyRequest: Response{
		Code: http.StatusTooManyRequests,
		Body: keyvalue.KeyValue{},
		Header: map[string]string{
			"Content-Type": "application/json",
		},
	},
}

type ResponseMapperCtx struct {
	Responses   map[string]Response
	Logger      logger.Logger
	DefaultCode *DefaultResponseCode
}

func (m *ResponseMapperCtx) Load(rs map[string]Response) {
	for key, val := range rs {
		m.Responses[key] = val
	}
}

func (m *ResponseMapperCtx) GetSuccess() Response {
	if m.DefaultCode == nil {
		return builtInResponse.ResponseSuccess
	}

	if m.DefaultCode.Success == "" {
		return builtInResponse.ResponseSuccess
	}

	r, exist := m.Responses[m.DefaultCode.Success]
	if !exist {
		m.Logger.Debug(fmt.Sprintf("Response Code: %s not mapped", m.DefaultCode.Success), nil)
		return builtInResponse.ResponseSuccess
	}

	return r
}

func (m *ResponseMapperCtx) GetInternalError() Response {
	if m.DefaultCode == nil {
		return builtInResponse.ResponseInternalError
	}

	if m.DefaultCode.InternalError == "" {
		return builtInResponse.ResponseInternalError
	}

	r, exist := m.Responses[m.DefaultCode.InternalError]
	if !exist {
		m.Logger.Debug(fmt.Sprintf("Response Code: %s not mapped", m.DefaultCode.InternalError), nil)
		return builtInResponse.ResponseInternalError
	}

	return r
}

func (m *ResponseMapperCtx) Get(code string, options *struct{ Success bool }) Response {
	r, exist := m.Responses[code]

	if !exist {
		m.Logger.Debug(fmt.Sprintf("Response Code: %s not mapped", code), nil)
		if options != nil && options.Success {
			return m.GetSuccess()
		}

		return m.GetInternalError()
	}
	return r
}

func ResponseMapper(cfg ResponseMapperCfg) *ResponseMapperCtx {
	if cfg.Logger == nil {
		apperr.Throw(apperr.New("ResponseMapper Logger cannot be nil!"))
	}

	cfg.Logger.Debug("OK", nil)

	responses := map[string]Response{}
	builtInResponseVal := reflect.ValueOf(builtInResponse)

	for i := 0; i < builtInResponseVal.NumField(); i++ {
		r := builtInResponseVal.Field(i).Interface().(Response)
		code, _ := parser.Int(r.Code).ToString()
		responses[code] = r
	}

	m := &ResponseMapperCtx{
		Responses:   responses,
		Logger:      cfg.Logger,
		DefaultCode: cfg.DefaultCode,
	}

	return m
}
