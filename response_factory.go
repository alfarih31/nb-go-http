package noob

func NewResponseSuccess(body ResponseBody, header ...ResponseHeader) Response {
	return NewResponse(StatusOK, body, header...)
}

func NewResponseError(code HTTPStatusCode, body ResponseBody, header ...ResponseHeader) ResponseError {
	h := ResponseHeader{}
	if len(header) > 0 {
		h = header[0]
	}

	return &responseError{
		Response: &response{
			Code:   &code,
			Body:   &body,
			Header: &h,
		},
	}
}
