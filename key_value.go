package nbgohttp

import (
	"encoding/json"
)

type KeyValue map[string]interface{}

func (k KeyValue) AssignTo(target KeyValue) {
	for key, val := range k {
		target[key] = val
	}
}

func (k KeyValue) Assign(source KeyValue) {
	for key, val := range source {
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
		ThrowError(Err{
			Err: e,
		})
	}

	var t map[string]interface{}

	e = json.Unmarshal(j, &t)

	if e != nil {
		ThrowError(Err{
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
