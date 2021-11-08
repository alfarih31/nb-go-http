package nbgohttp

import (
	"encoding/json"
)

type ResponseHeader = KeyValue

type ResponseStatus struct {
	Code          int    `json:"code"`
	MessageClient string `json:"message_client"`
	MessageServer string `json:"message_server"`
}

type ResponseBody struct {
	Status ResponseStatus `json:"status"`
	Meta   interface{}    `json:"meta,omitempty"`
	Data   interface{}    `json:"data,omitempty"`
	Errors []interface{}  `json:"errors,omitempty"`
}

type Response struct {
	Code   int
	Header ResponseHeader
	Body   ResponseBody
}

func (s ResponseStatus) AssignTo(target *ResponseStatus) {
	if target.MessageClient == "" {
		target.MessageClient = s.MessageClient
	}

	if target.MessageServer == "" {
		target.MessageServer = s.MessageServer
	}

	if target.Code == 0 {
		target.Code = s.Code
	}
}

func (r Response) Compose(res Response) Response {
	r.Header.AssignTo(res.Header)
	res.Body = r.Body.Compose(res.Body)

	if res.Code == 0 {
		res.Code = r.Code
	}

	return Response{
		Code:   res.Code,
		Header: res.Header,
		Body:   res.Body,
	}
}

func (b ResponseBody) Compose(body ResponseBody) ResponseBody {

	b.Status.AssignTo(&body.Status)

	if body.Errors == nil {
		body.Errors = b.Errors
	}

	return body
}

func (b ResponseBody) String() string {
	j, err := json.Marshal(b)

	if err != nil {
		return ""
	}

	return string(j)
}

const (
	StatusSuccess           = "200"
	StatusErrorBadRequest   = "400"
	StatusErrorUnauthorized = "401"
	StatusErrorForbidden    = "403"
	StatusErrorNotFound     = "404"
	StatusErrorInternal     = "500"
	StatusErrorBadGateway   = "502"
)
