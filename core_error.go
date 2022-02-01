package noob

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DataDog/gostackparse"
	"github.com/alfarih31/nb-go-keyvalue"
	"reflect"
	"runtime"
	"runtime/debug"
	"strings"
)

const DefaultErrCode = 1

type CoreError struct {
	Err    error                     `json:"err"`
	Code   uint                      `json:"code"`
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

func (msg *CoreError) Trace() {
	msg.Stack = StackTrace()
}

func (msg CoreError) StackTrace() *runtime.Frames {
	if msg.Frames != nil {
		return msg.Frames
	}

	return GetRuntimeFrames(3)
}

func (msg *CoreError) Error() string {
	return msg.Err.Error()
}

func (msg *CoreError) Errors() interface{} {
	return msg.Err
}

func (msg CoreError) Compose(message string, data ...interface{}) *CoreError {
	e := msg
	e.Stack = StackTrace()

	if message != "" {
		e.Err = errors.New(message)
	}

	e.Meta = data

	return &e
}

func NewCoreError(msg string, meta ...interface{}) *CoreError {
	e := &CoreError{
		Err:  errors.New(msg),
		Code: DefaultErrCode,
		Meta: meta,
	}

	return e
}

func (msg *CoreError) JSON() interface{} {
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

func (msg CoreError) MarshalJSON() ([]byte, error) {
	return json.Marshal(msg.JSON())
}

type Errors []*CoreError

func (e Errors) MarshalJSON() ([]byte, error) {
	jsonData := make([]interface{}, len(e))
	for i, er := range e {
		jsonData[i] = er.JSON()
	}

	return json.Marshal(jsonData)
}

func (e Errors) String() string {
	if len(e) == 0 {
		return ""
	}
	var buffer strings.Builder
	for i, msg := range e {
		fmt.Fprintf(&buffer, "Error #%02d: %v\n", i+1, msg)
	}
	return buffer.String()
}
