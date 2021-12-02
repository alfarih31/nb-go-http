package apperr

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/DataDog/gostackparse"
	"github.com/alfarih31/nb-go-http/keyvalue"
	"reflect"
	"runtime"
	"runtime/debug"
)

const DefaultErrCode = "_"

type AppErr struct {
	Err    error                     `json:"err"`
	Code   string                    `json:"code"`
	Meta   interface{}               `json:"meta"`
	Stack  []*gostackparse.Goroutine `json:"_stacks,omitempty"`
	Frames *runtime.Frames           `json:"_frames,omitempty"`
}

func GetRuntimeFrames(skip int) *runtime.Frames {
	pc := make([]uintptr, 100, 100)

	if skip == 0 {
		skip = 2
	}

	n := runtime.Callers(skip, pc)
	if n == 0 {
		return &runtime.Frames{}
	}

	return runtime.CallersFrames(pc[:n])
}

func StackTrace() []*gostackparse.Goroutine {
	stacks, _ := gostackparse.Parse(bytes.NewReader(debug.Stack()))

	return stacks
}

func (msg *AppErr) Trace() {
	msg.Stack = StackTrace()
}

func (msg AppErr) StackTrace() *runtime.Frames {
	if msg.Frames != nil {
		return msg.Frames
	}

	return GetRuntimeFrames(3)
}

func (msg *AppErr) Error() string {
	return msg.Err.Error()
}

func (msg *AppErr) Errors() interface{} {
	return msg.Err
}

func (msg AppErr) Throw(message string, data ...interface{}) {
	e := msg.Compose(message, data)

	panic(&e)
}

func (msg AppErr) Compose(message string, data ...interface{}) AppErr {
	e := msg
	e.Stack = StackTrace()

	if message != "" {
		e.Err = errors.New(message)
	}

	e.Meta = data

	return msg
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
