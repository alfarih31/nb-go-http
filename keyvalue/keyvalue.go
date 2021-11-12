package keyvalue

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

type KeyValue map[string]interface{}

func (k KeyValue) AssignTo(target KeyValue, replaceExist bool) {
	for key, val := range k {
		targetValue, exist := target[key]

		// Recursive assignTo
		if reflect.ValueOf(val).Kind() == reflect.Map && reflect.ValueOf(targetValue).Kind() == reflect.Map {
			sourceKvVal, _ := KeyValueFromStruct(val)
			targetKvVal, _ := KeyValueFromStruct(targetValue)

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
			sourceKvVal, _ := KeyValueFromStruct(val)
			existingKvVal, _ := KeyValueFromStruct(existingValue)

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

func StructToMap(strct interface{}) (map[string]interface{}, error) {
	tStruct := reflect.TypeOf(strct).Kind()
	if tStruct == reflect.Map {
		return strct.(map[string]interface{}), nil
	}

	if tStruct != reflect.Struct {
		return nil, errors.New(fmt.Sprintf("val not a struct type: %s", tStruct))
	}

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

func KeyValueFromStruct(strct interface{}) (KeyValue, error) {
	mapString, err := StructToMap(strct)

	if err != nil {
		return nil, err
	}

	kv := KeyValue{}
	for key, val := range mapString {
		kv[key] = val
	}

	return kv, nil
}
