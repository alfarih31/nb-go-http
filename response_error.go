package noob

var _ ResponseError = new(responseError)

type ResponseError interface {
	Response
	Error() string
}

type responseError struct {
	Response
}

func (e *responseError) Error() string {
	body := e.GetBody()
	return body.Message
}
