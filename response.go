package noob

import (
	"github.com/alfarih31/nb-go-keyvalue"
)

type Response interface {
	Compose(sourceRes Response, replaceExist ...bool) Response
	ComposeBody(body interface{}, replaceExist ...bool) Response
	GetBody() interface{}
	SetBody(b interface{}) Response
	GetHeader() map[string]string
	SetHeader(h map[string]string) Response
	GetCode() uint
	SetCode(c uint) Response
	Copy() Response
}

type DefaultResponse struct {
	Code   uint
	Header map[string]string
	Body   interface{}
}

func (r *DefaultResponse) GetBody() interface{} {
	return r.Body
}

func (r *DefaultResponse) GetHeader() map[string]string {
	if r.Header == nil {
		return nil
	}

	return r.Header
}

func (r *DefaultResponse) GetCode() uint {
	return r.Code
}

func (r *DefaultResponse) SetBody(b interface{}) Response {
	r.Body = b

	return r
}

func (r *DefaultResponse) SetHeader(h map[string]string) Response {
	r.Header = h
	return r
}

func (r *DefaultResponse) SetCode(c uint) Response {
	r.Code = c

	return r
}

func (r *DefaultResponse) Compose(sourceRes Response, replaceExist ...bool) Response {
	rExist := false
	if len(replaceExist) > 0 {
		rExist = replaceExist[0]
	}

	if sourceRes.GetHeader() != nil {
		if !rExist && r.GetHeader() == nil {
			r.SetHeader(sourceRes.GetHeader())
		} else {
			r.SetHeader(sourceRes.GetHeader())
		}
	}

	r.ComposeBody(sourceRes.GetBody(), rExist)

	if sourceRes.GetCode() != 0 {
		if !rExist && r.GetCode() == 0 {
			r.SetCode(sourceRes.GetCode())
		} else {
			r.SetCode(sourceRes.GetCode())
		}
	}

	return r
}

func (r *DefaultResponse) Copy() Response {
	return &DefaultResponse{
		Code:   r.GetCode(),
		Header: r.GetHeader(),
		Body:   r.GetBody(),
	}
}

func (r *DefaultResponse) ComposeBody(body interface{}, replaceExist ...bool) Response {
	rExist := false
	if len(replaceExist) > 0 {
		rExist = replaceExist[0]
	}

	if body == nil {
		return r
	}

	tBody := r.GetBody()
	if tBody == nil {
		return r
	}

	sourceBody, err := keyvalue.FromStruct(body)
	if err != nil {
		return r
	}

	targetBody, err := keyvalue.FromStruct(tBody)
	if err != nil {
		return r
	}

	targetBody.Assign(sourceBody, rExist)

	err = targetBody.Unmarshal(&tBody)
	if err != nil {
		return r
	}
	return r.SetBody(tBody)
}
