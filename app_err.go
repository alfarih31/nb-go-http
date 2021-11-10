package nbgohttp

import (
	"bytes"
	"github.com/DataDog/gostackparse"
	"runtime/debug"
)

const AppErrName = "AppError"

type Err struct {
	Name    string                    `json:"name"`
	Message string                    `json:"message"`
	Code    string                    `json:"code"`
	Data    interface{}               `json:"data,omitempty"`
	Err     interface{}               `json:"errors,omitempty"`
	Stack   []*gostackparse.Goroutine `json:"_stack,omitempty"`
}

func StackTrace() []*gostackparse.Goroutine {
	stacks, _ := gostackparse.Parse(bytes.NewReader(debug.Stack()))

	return stacks
}

func (e Err) Error() string {
	return e.Message
}

func (e Err) Errors() interface{} {
	return e.Err
}

func (e Err) Throw(er *Err) {
	if er == nil {
		er = &e
	}

	er.Stack = StackTrace()
	if er.Name == "" {
		er.Name = e.Name
	}

	if er.Code == "" {
		er.Code = e.Code
	}

	if er.Message == "" {
		er.Message = e.Message
	}

	if er.Data == nil {
		er.Data = e.Data
	}

	if er.Err == nil {
		er.Err = e.Err
	}

	panic(er)
}

func ThrowError(e *Err) {
	e.Stack = StackTrace()
	if e.Name == "" {
		e.Name = AppErrName
	}

	panic(e)
}

func NewError(e Err) *Err {
	er := &Err{
		Code:    e.Code,
		Name:    e.Name,
		Message: e.Message,
		Data:    e.Data,
	}

	if er.Name == "" {
		er.Name = AppErrName
	}

	return er
}
