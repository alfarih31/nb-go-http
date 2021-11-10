package nbgohttp

import (
	"context"
	"fmt"
	"net/http"
)

type ResponseMapperCfg struct {
	Context           context.Context
	Logger            ILogger
	SuccessCode       string
	InternalErrorCode string
}

type DefaultResponse struct {
	Success       Response
	InternalError Response
}

type ResponseMapperCtx struct {
	Context           context.Context
	Responses         map[string]Response
	Logger            ILogger
	successCode       string
	internalErrorCode string
	defaults          DefaultResponse
}

func (m *ResponseMapperCtx) Load(rs map[string]Response) {
	for key, val := range rs {
		m.Responses[key] = val
	}
}

func (m *ResponseMapperCtx) GetSuccess() Response {
	if m.successCode == "" {
		m.Logger.Debug("successCode is ''", nil)
		return m.defaults.Success
	}

	r, exist := m.Responses[m.successCode]
	if !exist {
		m.Logger.Debug(fmt.Sprintf("Response Code: %s not mapped", m.successCode), nil)
		return m.defaults.Success
	}

	return r
}

func (m *ResponseMapperCtx) GetInternalError() Response {
	if m.internalErrorCode == "" {
		m.Logger.Debug("internalErrorCode is ''", nil)
		return m.defaults.InternalError
	}

	r, exist := m.Responses[m.internalErrorCode]
	if !exist {
		m.Logger.Debug(fmt.Sprintf("Response Code: %s not mapped", m.internalErrorCode), nil)
		return m.defaults.InternalError
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

func (m *ResponseMapperCtx) WithContext(ctx context.Context) *ResponseMapperCtx {
	rm := ResponseMapper(ResponseMapperCfg{
		Context:           ctx,
		Logger:            m.Logger,
		SuccessCode:       m.successCode,
		InternalErrorCode: m.internalErrorCode,
	})
	rm.Responses = m.Responses
	return m
}

func ResponseMapper(cfg ResponseMapperCfg) *ResponseMapperCtx {
	if cfg.Logger == nil {
		ThrowError(&Err{Message: "ResponseMapper Logger cannot be nil!"})
	}

	cfg.Logger.Debug("OK", nil)

	m := &ResponseMapperCtx{
		Context:           cfg.Context,
		Responses:         map[string]Response{},
		Logger:            cfg.Logger,
		successCode:       cfg.SuccessCode,
		internalErrorCode: cfg.InternalErrorCode,
		defaults: DefaultResponse{
			Success: Response{
				Code: http.StatusOK,
			},
			InternalError: Response{
				Code: http.StatusInternalServerError,
			},
		},
	}

	return m
}
