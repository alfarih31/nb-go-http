package noob

import (
	"fmt"
	"github.com/alfarih31/nb-go-http/parser"
)

type QueryParser HandlerCtx

type QueryParserOption struct {
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

func (p QueryParser) GetString(key string, opt ...QueryParserOption) (string, error) {
	val := p.Query(key)

	if val == "" {
		optVal, err := getOptions(key, opt)
		if err != nil {
			return "", err
		}

		if optVal == nil {
			return "", nil
		}

		return optVal.(string), nil
	}

	return val, nil
}

func (p QueryParser) GetInt(key string, opt ...QueryParserOption) (int, error) {
	val := p.Query(key)

	if val == "" {
		optVal, err := getOptions(key, opt)
		if err != nil {
			return 0, err
		}

		if optVal == nil {
			return 0, nil
		}

		return optVal.(int), nil
	}

	i, err := parser.String(val).ToInt()

	if err != nil {
		optVal, e := getOptions(key, opt)
		if e != nil {
			return 0, e
		}

		if optVal == nil {
			return 0, nil
		}

		return optVal.(int), nil
	}

	return i, err
}

func (p QueryParser) GetBool(key string, opt ...QueryParserOption) (bool, error) {
	val := p.Query(key)

	if val == "" {
		optVal, err := getOptions(key, opt)
		if err != nil {
			return false, err
		}

		if optVal == nil {
			return false, nil
		}

		return optVal.(bool), nil
	}

	b, err := parser.String(val).ToBool()

	if err != nil {
		optVal, e := getOptions(key, opt)
		if e != nil {
			return false, e
		}

		if optVal == nil {
			return false, nil
		}

		return optVal.(bool), nil
	}

	return b, err
}
