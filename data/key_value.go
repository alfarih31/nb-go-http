package data

import (
	"encoding/json"
	"github.com/alfarih31/nb-go-http/app_error"
	"reflect"
)

type KeyValue map[string]interface{}

func (k KeyValue) AssignTo(target KeyValue, replaceExist bool) {
	for key, val := range k {
		targetValue, exist := target[key]

		// Recursive assignTo
		if reflect.ValueOf(val).Kind() == reflect.Map && reflect.ValueOf(targetValue).Kind() == reflect.Map {
			sourceKvVal := KeyValueFromStruct(val)
			targetKvVal := KeyValueFromStruct(targetValue)

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
			sourceKvVal := KeyValueFromStruct(val)
			existingKvVal := KeyValueFromStruct(existingValue)

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

func StructToMapString(strct interface{}) map[string]interface{} {
	j, e := json.Marshal(strct)

	if e != nil {
		apperror.ThrowError(&apperror.Err{
			Err: e,
		})
	}

	var t map[string]interface{}

	e = json.Unmarshal(j, &t)

	if e != nil {
		apperror.ThrowError(&apperror.Err{
			Err: e,
		})
	}

	return t
}

func KeyValueFromStruct(strct interface{}) KeyValue {
	kv := KeyValue{}
	for key, val := range StructToMapString(strct) {
		kv[key] = val
	}

	return kv
}
