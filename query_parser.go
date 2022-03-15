package noob

import (
	"fmt"
	keyvalue "github.com/alfarih31/nb-go-keyvalue"
	"github.com/alfarih31/nb-go-parser"
)

type QueryParser HandlerCtx

type QueryParserOption struct {
	Default  interface{}
	Required bool
}

type QueryValueType uint

const (
	QueryValueTypeString QueryValueType = iota
	QueryValueTypeBool
	QueryValueTypeInt
	QueryValueTypeInt32
	QueryValueTypeInt64
)

type Query struct {
	Key      string
	Type     QueryValueType
	Default  interface{}
	Required bool
}

func getOptions(key string, opt []QueryParserOption) (interface{}, error) {
	if len(opt) > 0 {
		o := opt[0]
		if o.Required {
			if o.Default != nil {
				return o.Default, nil
			}

			return nil, fmt.Errorf("qs: %s is required", key)
		}
	}

	return nil, nil
}

func getKeyErr(key string, err error) error {
	return fmt.Errorf("qs: %s error, %v", key, err)
}

func (p QueryParser) GetQueries(target interface{}, qs []Query) error {
	kv := keyvalue.KeyValue{}
	for _, q := range qs {
		var (
			v   interface{}
			err error
		)
		switch q.Type {
		case QueryValueTypeString:
			v, err = p.GetString(q.Key, QueryParserOption{Default: q.Default, Required: q.Required})
		case QueryValueTypeBool:
			v, err = p.GetBool(q.Key, QueryParserOption{Default: q.Default, Required: q.Required})
		case QueryValueTypeInt:
			v, err = p.GetInt(q.Key, QueryParserOption{Default: q.Default, Required: q.Required})
		case QueryValueTypeInt32:
			v, err = p.GetInt32(q.Key, QueryParserOption{Default: q.Default, Required: q.Required})
		case QueryValueTypeInt64:
			v, err = p.GetInt64(q.Key, QueryParserOption{Default: q.Default, Required: q.Required})
		default:
			Log.Warn(fmt.Sprintf("Unknown QueryValueType, Key=%s, Type=%d", q.Key, q.Type))
		}

		if err != nil {
			return err
		}

		kv[q.Key] = v
	}

	return kv.Unmarshal(target)
}

func (p QueryParser) GetString(key string, opt ...QueryParserOption) (*string, error) {
	val := p.Query(key)

	if val == "" {
		optVal, err := getOptions(key, opt)
		if err != nil {
			return nil, err
		}

		if optVal != nil {
			v := optVal.(string)
			return &v, nil
		}

		return nil, nil
	}

	return &val, nil
}

func (p QueryParser) GetInt(key string, opt ...QueryParserOption) (*int, error) {
	val := p.Query(key)

	optVal, optErr := getOptions(key, opt)

	if val == "" {
		if optErr != nil {
			return nil, optErr
		}

		if optVal != nil {
			v := optVal.(int)

			return &v, nil
		}

		return nil, nil
	}

	i, err := parser.String(val).ToInt()

	if err != nil {
		if optErr != nil {
			return nil, getKeyErr(key, err)
		}

		if optVal != nil {
			v := optVal.(int)

			return &v, nil
		}

		return nil, getKeyErr(key, err)
	}

	return &i, err
}

func (p QueryParser) GetInt32(key string, opt ...QueryParserOption) (*int32, error) {
	val := p.Query(key)

	optVal, optErr := getOptions(key, opt)

	if val == "" {
		if optErr != nil {
			return nil, optErr
		}

		if optVal != nil {
			v := optVal.(int32)

			return &v, nil
		}

		return nil, nil
	}

	i, err := parser.String(val).ToInt()
	i32 := int32(i)

	if err != nil {
		if optErr != nil {
			return nil, getKeyErr(key, err)
		}

		if optVal != nil {
			v := optVal.(int32)

			return &v, nil
		}

		return nil, getKeyErr(key, err)
	}

	return &i32, err
}

func (p QueryParser) GetInt64(key string, opt ...QueryParserOption) (*int64, error) {
	val := p.Query(key)

	optVal, optErr := getOptions(key, opt)

	if val == "" {
		if optErr != nil {
			return nil, optErr
		}

		if optVal != nil {
			v := optVal.(int64)

			return &v, nil
		}

		return nil, nil
	}

	i, err := parser.String(val).ToInt()
	i64 := int64(i)

	if err != nil {
		if optErr != nil {
			return nil, getKeyErr(key, err)
		}

		if optVal != nil {
			v := optVal.(int64)

			return &v, nil
		}

		return nil, getKeyErr(key, err)
	}

	return &i64, err
}

func (p QueryParser) GetBool(key string, opt ...QueryParserOption) (*bool, error) {
	val := p.Query(key)

	optVal, optErr := getOptions(key, opt)

	if val == "" {
		if optErr != nil {
			return nil, optErr
		}

		if optVal != nil {
			v := optVal.(bool)

			return &v, nil
		}

		return nil, nil
	}

	b, err := parser.String(val).ToBool()

	if err != nil {
		if optErr != nil {
			return nil, getKeyErr(key, err)
		}

		if optVal != nil {
			v := optVal.(bool)

			return &v, nil
		}

		return nil, getKeyErr(key, err)
	}

	return &b, err
}
