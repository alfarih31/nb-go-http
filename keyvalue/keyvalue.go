package keyvalue

import (
	"encoding/json"
	"errors"
	"reflect"
)

type KeyValue map[string]interface{}

func (k KeyValue) AssignTo(target KeyValue, replaceExist bool) {
	for key, val := range k {
		targetValue, exist := target[key]

		// Recursive assignTo
		if reflect.ValueOf(val).Kind() == reflect.Map && reflect.ValueOf(targetValue).Kind() == reflect.Map {
			sourceKvVal, _ := FromStruct(val)
			targetKvVal, _ := FromStruct(targetValue)

			sourceKvVal.AssignTo(targetKvVal, replaceExist)

			target[key] = targetKvVal
			return
		}

		// Check is target is zero value
		t := reflect.TypeOf(targetValue)
		isZero := true
		if t != nil {
			isZero = targetValue == reflect.Zero(t).Interface()
		}

		// If exist but don't replace exist then continue
		if exist && !replaceExist && !isZero {
			continue
		}

		target[key] = val
	}
}

func (k KeyValue) Assign(source KeyValue, replaceExist bool) {
	for key, val := range source {
		existingValue, exist := k[key]

		// Recursive assign
		if reflect.ValueOf(val).Kind() == reflect.Map && reflect.ValueOf(existingValue).Kind() == reflect.Map {
			sourceKvVal, _ := FromStruct(val)
			existingKvVal, _ := FromStruct(existingValue)

			existingKvVal.Assign(sourceKvVal, replaceExist)

			k[key] = existingKvVal
			return
		}

		// Check existing is zero value
		t := reflect.TypeOf(existingValue)
		isZero := true
		if t != nil {
			isZero = existingValue == reflect.Zero(t).Interface()
		}

		// If exist & not zero value but don't replace exist then continue
		if exist && !replaceExist && !isZero {
			continue
		}

		k[key] = val
	}
}

func (k KeyValue) Keys() []string {
	var keys []string
	for key := range k {
		keys = append(keys, key)
	}
	return keys
}

func (k KeyValue) Values() []interface{} {
	var values []interface{}
	for _, val := range k {
		values = append(values, val)
	}

	return values
}

func (k KeyValue) String() string {
	j, _ := json.Marshal(k)
	return string(j)
}

func IsAbleToConvert(p interface{}) bool {
	t := reflect.TypeOf(p)
	name := t.Name()
	kind := t.Kind()

	if name == "KeyValue" {
		return true
	}

	switch kind {
	case reflect.Map:
		fallthrough
	case reflect.Struct:
		return true
	}

	return true
}

func structToMap(strct interface{}) (map[string]interface{}, error) {
	j, e := json.Marshal(strct)

	if e != nil {
		return nil, e
	}

	var t map[string]interface{}

	e = json.Unmarshal(j, &t)

	if e != nil {
		return nil, e
	}

	return t, nil
}

func FromStruct(strct interface{}) (KeyValue, error) {
	if !IsAbleToConvert(strct) {
		return nil, errors.New("cannot convert")
	}

	mapString, err := structToMap(strct)

	if err != nil {
		return nil, err
	}

	kv := KeyValue{}
	for key, val := range mapString {
		kv[key] = val
	}

	return kv, nil
}
