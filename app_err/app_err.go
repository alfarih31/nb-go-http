package apperr

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/DataDog/gostackparse"
	"github.com/alfarih31/nb-go-http/keyvalue"
	"reflect"
	"runtime/debug"
)

const DefaultErrCode = "_"

type AppErr struct {
	Err   error                     `json:"err"`
	Code  string                    `json:"code"`
	Meta  interface{}               `json:"meta"`
	Stack []*gostackparse.Goroutine `json:"_stack,omitempty"`
}

func StackTrace() []*gostackparse.Goroutine {
	stacks, _ := gostackparse.Parse(bytes.NewReader(debug.Stack()))

	return stacks
}

func (msg *AppErr) Trace() {
	msg.Stack = StackTrace()
}

func (msg *AppErr) Error() string {
	return msg.Err.Error()
}

func (msg *AppErr) Errors() interface{} {
	return msg.Err
}

func (msg AppErr) Throw(message string, data ...interface{}) {
	msg.Stack = StackTrace()

	if message != "" {
		msg.Err = errors.New(message)
	}

	msg.Meta = data

	panic(&msg)
}

func Throw(e *AppErr) {
	e.Stack = StackTrace()

	panic(e)
}

func New(msg string, meta ...interface{}) *AppErr {
	e := &AppErr{
		Err:  errors.New(msg),
		Code: DefaultErrCode,
		Meta: meta,
	}

	return e
}

func (msg *AppErr) JSON() interface{} {
	jsonData := keyvalue.KeyValue{}
	if msg.Meta != nil {
		value := reflect.ValueOf(msg.Meta)
		switch value.Kind() {
		case reflect.Struct:
			return msg.Meta
		case reflect.Map:
			for _, key := range value.MapKeys() {
				jsonData[key.String()] = value.MapIndex(key).Interface()
			}
		default:
			jsonData["meta"] = msg.Meta
		}
	}
	if _, ok := jsonData["message"]; !ok {
		jsonData["message"] = msg.Error()
	}

	jsonData["_stack"] = msg.Stack

	return jsonData
}

func (msg AppErr) MarshalJSON() ([]byte, error) {
	return json.Marshal(msg.JSON())
}
