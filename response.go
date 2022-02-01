package noob

import (
	keyvalue "github.com/alfarih31/nb-go-keyvalue"
	parser "github.com/alfarih31/nb-go-parser"
)

type ResponseMap map[HTTPStatusCode]Response

type Response interface {
	Compose(sourceRes Response, replaceExist ...bool) Response
	ComposeBody(body ResponseBody, replaceExist ...bool) Response
	ComposeHeader(h ResponseHeader, replaceExist ...bool) Response
	GetBody() *ResponseBody
	GetHeader() *ResponseHeader
	GetCode() *HTTPStatusCode
	Copy() Response
}

type ResponseHeader map[string]string

func (r ResponseHeader) Copy() *ResponseHeader {
	nh := ResponseHeader{}
	for k, v := range r {
		nh[k] = v
	}

	return &nh
}

type ResponseBody struct {
	Code    uint        `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"_error,omitempty"`
}

func (b ResponseBody) Copy() *ResponseBody {
	return &ResponseBody{
		Code:    b.Code,
		Message: b.Message,
		Data:    b.Data,
		Errors:  b.Errors,
	}
}

type response struct {
	Code   *HTTPStatusCode
	Header *ResponseHeader
	Body   *ResponseBody
}

func (r *response) GetBody() *ResponseBody {
	return r.Body
}

func (r *response) GetHeader() *ResponseHeader {
	if r.Header == nil {
		return nil
	}

	return r.Header
}

func (r *response) GetCode() *HTTPStatusCode {
	return r.Code
}

func (r *response) Compose(sourceRes Response, replaceExist ...bool) Response {
	rExist := parser.GetOptBoolArg(replaceExist)

	r.ComposeBody(*sourceRes.GetBody(), rExist)

	r.ComposeHeader(*sourceRes.GetHeader(), rExist)

	if rExist {
		c := r.GetCode()
		*c = *sourceRes.GetCode()
	} else {
		if sourceRes.GetCode() != nil && r.GetCode() == nil {
			c := r.GetCode()
			*c = *sourceRes.GetCode()
		}
	}

	return r
}

func (r *response) ComposeHeader(h ResponseHeader, replaceExist ...bool) Response {
	rExist := parser.GetOptBoolArg(replaceExist)
	ch := *r.Header
	for k, v := range h {
		_, exist := ch[k]
		if exist && rExist {
			ch[k] = v
			continue
		}

		ch[k] = v
	}

	r.Header = &ch

	return r
}

func (r *response) Copy() Response {
	return &response{
		Code:   r.GetCode().Copy(),
		Header: r.GetHeader().Copy(),
		Body:   r.GetBody().Copy(),
	}
}

func (r *response) ComposeBody(body ResponseBody, replaceExist ...bool) Response {
	rExist := parser.GetOptBoolArg(replaceExist)

	tBody := r.GetBody()

	sourceBody, err := keyvalue.FromStruct(body)
	if err != nil {
		logR.Error(err)
		return r
	}

	targetBody, err := keyvalue.FromStruct(tBody)
	if err != nil {
		logR.Error(err)
		return r
	}

	targetBody.Assign(sourceBody, rExist)

	err = targetBody.Unmarshal(&tBody)
	if err != nil {
		logR.Error(err)
		return r
	}

	return r
}

func NewResponse(code HTTPStatusCode, body ResponseBody, header ...ResponseHeader) Response {
	h := ResponseHeader{}
	if len(header) > 0 {
		h = header[0]
	}

	return &response{
		Code:   &code,
		Body:   &body,
		Header: &h,
	}
}
