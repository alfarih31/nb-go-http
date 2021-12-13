package noob

import (
	"github.com/alfarih31/nb-go-keyvalue"
)

type Response struct {
	Code   int
	Header map[string]string
	Body   interface{}
}

func (r Response) ComposeTo(res *Response) {
	if res.Header == nil {
		res.Header = r.Header
	}

	res.ComposeBody(r.Body)

	if res.Code == 0 {
		res.Code = r.Code
	}
}

func (r Response) Compose(res Response) *Response {
	outRes := &Response{
		Code:   res.Code,
		Header: res.Header,
		Body:   res.Body,
	}

	if outRes.Header == nil {
		outRes.Header = r.Header
	}

	outRes.ComposeBody(r.Body)

	if outRes.Code == 0 {
		outRes.Code = r.Code
	}

	return outRes
}

func (r *Response) ComposeBody(body interface{}) {
	if body == nil {
		return
	}

	if r.Body == nil {
		return
	}

	sourceBody, err := keyvalue.FromStruct(body)
	if err != nil {
		r.Body = body
		return
	}

	targetBody, err := keyvalue.FromStruct(r.Body)
	if err != nil {
		return
	}

	targetBody.Assign(sourceBody, false)

	r.Body = targetBody
}
