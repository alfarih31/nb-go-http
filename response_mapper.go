package nbgohttp

import (
	"fmt"
	"net/http"
)

type TResponseMapper struct {
	Responses map[string]Response
	Logger    ILogger
}

type IResponseMapper interface {
	Load(rs map[string]Response)
	GetSuccess() Response
	GetInternalError() Response
	Get(code string, options *struct{ Success bool }) Response
}

func initStandardResponses() map[string]Response {
	return map[string]Response{
		StatusSuccess: {
			Code: http.StatusOK,
			Body: ResponseBody{
				Status: ResponseStatus{
					Code:          0,
					MessageClient: "Success",
					MessageServer: "Success",
				},
			},
		},
		StatusErrorBadRequest: {
			Code: http.StatusBadRequest,
			Body: ResponseBody{
				Status: ResponseStatus{
					Code:          1,
					MessageClient: "Bad Request",
					MessageServer: "Bad Request",
				},
			},
		},
		StatusErrorUnauthorized: {
			Code: http.StatusUnauthorized,
			Body: ResponseBody{
				Status: ResponseStatus{
					Code:          1,
					MessageClient: "Unauthorized",
					MessageServer: "Unauthorized",
				},
			},
		},
		StatusErrorForbidden: {
			Code: http.StatusForbidden,
			Body: ResponseBody{
				Status: ResponseStatus{
					Code:          1,
					MessageClient: "Forbidden",
					MessageServer: "Forbidden",
				},
			},
		},
		StatusErrorNotFound: {
			Code: http.StatusNotFound,
			Body: ResponseBody{
				Status: ResponseStatus{
					Code:          1,
					MessageClient: "Not Found",
					MessageServer: "Not Found",
				},
			},
		},
		StatusErrorInternal: {
			Code: http.StatusInternalServerError,
			Body: ResponseBody{
				Status: ResponseStatus{
					Code:          1,
					MessageClient: "Internal Error",
					MessageServer: "Internal Error",
				},
			},
		},
		StatusErrorBadGateway: {
			Code: http.StatusBadGateway,
			Body: ResponseBody{
				Status: ResponseStatus{
					Code:          1,
					MessageClient: "Bad Gateway",
					MessageServer: "Bad Gateway",
				},
			},
		},
	}
}

func (m *TResponseMapper) Load(rs map[string]Response) {
	for key, val := range rs {
		m.Responses[key] = val
	}
}

func (m *TResponseMapper) GetSuccess() Response {
	return m.Responses[StatusSuccess]
}

func (m *TResponseMapper) GetInternalError() Response {
	return m.Responses[StatusErrorInternal]
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

func ResponseMapper(logger ILogger) IResponseMapper {
	logger.Debug("OK", nil)

	m := TResponseMapper{
		Responses: initStandardResponses(),
		Logger:    logger,
	}

	return &m
}
