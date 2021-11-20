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
			return nil, err
		}

		if optVal != nil {
			v := optVal.(int)

			return &v, nil
		}

		return nil, nil
	}

	return &i, err
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
			return nil, err
		}

		if optVal != nil {
			v := optVal.(bool)

			return &v, nil
		}

		return nil, nil
	}

	return &b, err
}
