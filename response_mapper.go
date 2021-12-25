package noob

import (
	"fmt"
	"github.com/alfarih31/nb-go-http/app_err"
	"github.com/alfarih31/nb-go-http/logger"
	"github.com/alfarih31/nb-go-keyvalue"
	"net/http"
)

type ResponseMapperCfg struct {
	Logger            logger.Logger
	DefaultHeader     *keyvalue.KeyValue
	StandardResponses ResponseMap
}

// Map of http std codes to response codes
type standardCodes map[uint]uint

type ResponseMap map[uint]Response

type ResponseMapperCtx struct {
	Responses     ResponseMap
	Logger        logger.Logger
	DefaultHeader *keyvalue.KeyValue
	standardCodes standardCodes
}

var builtInResponse = ResponseMap{
	http.StatusOK: &DefaultResponse{
		Code: http.StatusOK,
		Body: keyvalue.KeyValue{},
		Header: map[string]string{
			"Content-Type": "application/json",
		},
	},
	http.StatusInternalServerError: &DefaultResponse{
		Code: http.StatusInternalServerError,
		Body: keyvalue.KeyValue{},
		Header: map[string]string{
			"Content-Type": "application/json",
		},
	},
	http.StatusNotFound: &DefaultResponse{
		Code: http.StatusNotFound,
		Body: keyvalue.KeyValue{},
		Header: map[string]string{
			"Content-Type": "application/json",
		},
	},
	http.StatusTooManyRequests: &DefaultResponse{
		Code: http.StatusTooManyRequests,
		Body: keyvalue.KeyValue{},
		Header: map[string]string{
			"Content-Type": "application/json",
		},
	},
}

func (m *ResponseMapperCtx) AssignDefaultHeader(r Response) {
	if m.DefaultHeader == nil {
		return
	}

	h, err := keyvalue.FromStruct(r.GetHeader())
	if err == nil {
		m.DefaultHeader.AssignTo(h, false)

		r.SetHeader(h.JSON())
	}
}

func (m *ResponseMapperCtx) Load(rs ResponseMap) *ResponseMapperCtx {
	for key, val := range rs {
		m.Responses[key] = val
	}

	return m
}

func (m *ResponseMapperCtx) GetSuccess() Response {
	// Get from responses
	r, exist := m.Responses[m.standardCodes[http.StatusOK]]
	if !exist {
		return builtInResponse[http.StatusOK]
	}

	return r
}

func (m *ResponseMapperCtx) GetInternalError() Response {
	// Get from responses
	r, exist := m.Responses[m.standardCodes[http.StatusInternalServerError]]
	if !exist {
		// If not exist get from builtIn
		return builtInResponse[http.StatusInternalServerError]
	}

	return r
}

func (m *ResponseMapperCtx) Get(code uint, options ...struct{ Success bool }) Response {
	// Get from Responses
	var (
		r     Response
		exist bool
	)
	r, exist = m.Responses[code]
	if !exist {
		// If not exist get from builtIn
		r, exist = builtInResponse[code]
		if !exist {
			m.Logger.Debug(fmt.Sprintf("Response Code: %d not mapped", code), nil)

			if len(options) > 0 && options[0].Success {
				return m.GetSuccess()
			}

			return m.GetInternalError()
		}
	}

	return r
}

func ResponseMapper(cfg ResponseMapperCfg) *ResponseMapperCtx {
	if cfg.Logger == nil {
		apperr.Throw(apperr.New("ResponseMapper Logger cannot be nil!"))
	}

	cfg.Logger.Debug("OK", nil)

	// Init standard codes & responses
	responses := ResponseMap{}
	stdCodes := standardCodes{}
	if cfg.StandardResponses != nil {
		for k, v := range cfg.StandardResponses {
			stdCodes[v.GetCode()] = k
			responses[k] = v
		}
	}

	return &ResponseMapperCtx{
		Responses:     responses,
		Logger:        cfg.Logger,
		DefaultHeader: cfg.DefaultHeader,
		standardCodes: stdCodes,
	}
}
