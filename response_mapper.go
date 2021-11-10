package nbgohttp

import (
	"fmt"
	"net/http"
)

type ResponseMapperCfg struct {
	Logger            ILogger
	SuccessCode       string
	InternalErrorCode string
}

type DefaultResponse struct {
	Success       Response
	InternalError Response
}

type TResponseMapper struct {
	Responses         map[string]Response
	Logger            ILogger
	successCode       string
	internalErrorCode string
	defaults          DefaultResponse
}

type IResponseMapper interface {
	Load(rs map[string]Response)
	GetSuccess() Response
	GetInternalError() Response
	Get(code string, options *struct{ Success bool }) Response
}

func (m *TResponseMapper) Load(rs map[string]Response) {
	for key, val := range rs {
		m.Responses[key] = val
	}
}

func (m *TResponseMapper) GetSuccess() Response {
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

func (m *TResponseMapper) GetInternalError() Response {
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

func (m *TResponseMapper) Get(code string, options *struct{ Success bool }) Response {
	r := m.Responses[code]

	if &r == nil {
		m.Logger.Debug(fmt.Sprintf("Response Code: %s not mapped", code), nil)
		if options != nil && options.Success {
			return m.GetSuccess()
		}

		return m.GetInternalError()
	}
	return r
}

func ResponseMapper(cfg ResponseMapperCfg) IResponseMapper {
	if cfg.Logger == nil {
		ThrowError(&Err{Message: "ResponseMapper Logger cannot be nil!"})
	}

	cfg.Logger.Debug("OK", nil)

	m := TResponseMapper{
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

	return &m
}
