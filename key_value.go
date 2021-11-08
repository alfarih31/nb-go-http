package nbgohttp

import (
	"encoding/json"
)

type KeyValue map[string]interface{}

func (k KeyValue) AssignTo(tk KeyValue) {
	for key, val := range k {
		tk[key] = val
	}
}

func (k KeyValue) Assign(tk KeyValue) {
	for key, val := range tk {
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

//func (k KeyValue) MarshalJSON() ([]byte, error) {
//    o := map[string]interface{}{}
//    for key, val := range k {
//        o[key] = val
//    }
//
//    return json.Marshal(o)
//}

func StructToMapString(k interface{}) map[string]interface{} {
	j, e := json.Marshal(k)

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

func KeyValueFromStruct(o interface{}) KeyValue {
	kv := KeyValue{}
	for key, val := range StructToMapString(o) {
		kv[key] = val
	}

	return kv
}
