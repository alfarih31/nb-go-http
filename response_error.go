package noob

var _ ResponseError = new(responseError)

type ResponseError interface {
	Response
	CopyError() ResponseError
	SetMessage(msg string) ResponseError
	Error() string
}

type responseError struct {
	Response
}

func (e *responseError) Error() string {
	body := e.GetBody()
	return body.Message
}

func (e *responseError) CopyError() ResponseError {
	rc := e.Copy()
	return &responseError{
		Response: rc,
	}
}

func (e *responseError) SetMessage(msg string) ResponseError {
	if msg == "" {
		return e
	}

	ec := e.CopyError()

	ec.ComposeBody(ResponseBody{
		Message: msg,
	})

	return ec
}
