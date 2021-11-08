package nbgohttp

import (
	"encoding/json"
	"runtime"
)

const AppErrName = "AppError"

type Err struct {
	Name    string
	Message string
	Code    string
	Data    interface{}
	Err     interface{}
	stack   []Trace
}

type Trace struct {
	FuncName string `json:"func_name"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

func traceError() []Trace {
	tb := runtime.ReadTrace()

	st := runtime.Stack(tb, false)

	var t []Trace

	for i := 0; i < st; i++ {
		pc, file, line, ok := runtime.Caller(i)

		if ok {
			fr := runtime.Frame{
				PC:   pc,
				File: file,
				Line: line,
			}

			fp := runtime.FuncForPC(pc)

			t = append(t, Trace{File: fr.File, FuncName: fp.Name(), Line: fr.Line})
		}
	}

	return t[6 : len(t)-2]
}

func (e Err) Stack() []Trace {
	return traceError()
}

func (e Err) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"name":    e.Name,
		"code":    e.Code,
		"message": e.Message,
		"data":    e.Data,
		"errors":  e.Err,
		"stack":   e.Stack(),
	})
}

func (e Err) Error() string {
	return e.Message
}

func (e Err) Errors() interface{} {
	return e.Err
}

func (e *Err) Throw(er *Err) {
	if er != nil {
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

		panic(*er)
	}

	panic(*e)
}

func ThrowError(e Err) {
	if e.Name == "" {
		e.Name = AppErrName
	}

	panic(e)
}

func NewError(e Err) Err {
	if &e.Name == nil {
		e.Name = AppErrName
	}

	return e
}
