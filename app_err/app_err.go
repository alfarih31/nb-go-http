package apperr

import (
	"bytes"
	"github.com/DataDog/gostackparse"
	"runtime/debug"
)

const DefaultErrCode = "AppError"
const ErrName = "AppError"

type AppErr struct {
	Name    string                    `json:"name"`
	Message string                    `json:"message"`
	Code    string                    `json:"code"`
	Data    interface{}               `json:"data,omitempty"`
	AppErr  interface{}               `json:"errors,omitempty"`
	Stack   []*gostackparse.Goroutine `json:"_stack,omitempty"`
}

func StackTrace() []*gostackparse.Goroutine {
	stacks, _ := gostackparse.Parse(bytes.NewReader(debug.Stack()))

	return stacks
}

func (e AppErr) Error() string {
	return e.Message
}

func (e AppErr) Errors() interface{} {
	return e.AppErr
}

func (e AppErr) Throw(message string, data ...interface{}) {
	e.Stack = StackTrace()

	if message != "" {
		e.Message = message
	}

	d := make([]interface{}, len(data))
	copy(d, data[:])
	e.Data = d

	panic(e)
}

func Throw(e AppErr) {
	e.Stack = StackTrace()

	panic(e)
}

func New(msgCode ...string) AppErr {
	message := msgCode[0]
	code := msgCode[1]

	e := AppErr{
		Code:    code,
		Name:    ErrName,
		Message: message,
	}

	if e.Code == "" {
		e.Code = DefaultErrCode
	}

	return e
}
