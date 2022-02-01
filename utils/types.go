package utils

import "reflect"

func IsTypeEqual(v1 interface{}, v2 interface{}) bool {
	t1 := reflect.TypeOf(v1)
	t2 := reflect.TypeOf(v2)
	if t1 == nil {
		if t2 == nil {
			return true
		}

		return false
	}

	if t1.Name() != t2.Name() {
		return false
	}

	if t1.Kind() != t2.Kind() {
		return false
	}

	return true
}

func IsInterfaceNil(i interface{}) bool {
	return i == nil || (reflect.ValueOf(i).Kind() == reflect.Ptr && reflect.ValueOf(i).IsNil())
}
