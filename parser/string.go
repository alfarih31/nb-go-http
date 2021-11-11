package parser

import (
	"errors"
	"strconv"
	"strings"
)

type String string

func (s String) ToInt() (int, error) {
	return strconv.Atoi(string(s))
}

func (s String) ToBool() (bool, error) {
	switch string(s) {
	case "true":
		fallthrough
	case "True":
		fallthrough
	case "TRUE":
		fallthrough
	case "1":
		return true, nil
	case "false":
		fallthrough
	case "False":
		fallthrough
	case "FALSE":
		fallthrough
	case "0":
		return false, nil
	default:
		return false, errors.New("string cannot be converted to bool")
	}
}

func (s String) ToStringArr() ([]string, error) {

	var ss []string
	for _, v := range strings.Split(string(s), ",") {
		ss = append(ss, strings.TrimSpace(v))
	}

	return ss, nil
}

func (s String) ToIntArr() ([]int, error) {

	var is []int
	for _, v := range strings.Split(string(s), ",") {
		i, e := strconv.Atoi(strings.TrimSpace(v))

		if e != nil {
			return nil, e
		}

		is = append(is, i)
	}

	return is, nil
}
