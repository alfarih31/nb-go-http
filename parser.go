package nbgohttp

import (
	"strconv"
	"strings"
)

type StringParser struct {
	str string
}

func (s StringParser) ToInt() (int, error) {
	return strconv.Atoi(s.str)
}

func (s StringParser) ToBool() (bool, error) {
	switch s.str {
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
		return false, &Err{Message: "string cannot be converted to bool"}
	}
}

func (s StringParser) ToStringArr() ([]string, error) {

	ss := []string{}
	for _, v := range strings.Split(s.str, ",") {
		ss = append(ss, strings.TrimSpace(v))
	}

	return ss, nil
}

func (s StringParser) ToIntArr() ([]int, error) {

	is := []int{}
	for _, v := range strings.Split(s.str, ",") {
		i, e := strconv.Atoi(strings.TrimSpace(v))

		if e != nil {
			return nil, e
		}

		is = append(is, i)
	}

	return is, nil
}

type BoolParser struct {
	b bool
}

func (b BoolParser) ToString() string {
	if b.b {
		return "true"
	}

	return "false"
}

func (b BoolParser) ToInt() int {
	if b.b {
		return 1
	}

	return 0
}