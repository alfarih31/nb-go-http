package env

import (
	"encoding/json"
	"errors"
	"github.com/alfarih31/nb-go-parser"
	"github.com/joho/godotenv"
	"os"
)

var ErrorVarNotExist = errors.New("variable not exist")

type env struct {
	envs    map[string]string
	useEnvs bool
}

type Env interface {
	GetInt(k string, def ...int) (int, error)
	GetString(k string, def ...string) (string, error)
	GetBool(k string, def ...bool) (bool, error)
	GetStringArr(k string, def ...[]string) ([]string, error)
	GetIntArr(k string, def ...[]int) ([]int, error)
	Dump() (string, error)
}

func (c env) GetInt(k string, def ...int) (int, error) {
	cfg, exist := c.get(k)

	if !exist {
		if len(def) == 0 {
			return 0, ErrorVarNotExist
		}
		return def[0], nil
	}

	i, e := parser.String(cfg).ToInt()

	if e != nil {
		if len(def) == 0 {
			return 0, e
		}

		return def[0], nil
	}

	return i, e
}

func (c env) GetString(k string, def ...string) (string, error) {
	cfg, exist := c.get(k)
	if !exist {
		if len(def) == 0 {
			return "", ErrorVarNotExist
		}
		return def[0], nil
	}

	return cfg, nil
}

func (c env) GetBool(k string, def ...bool) (bool, error) {
	cfg, exist := c.get(k)
	if !exist {
		if len(def) == 0 {
			return false, ErrorVarNotExist
		}
		return def[0], nil
	}

	b, e := parser.String(cfg).ToBool()
	if e != nil {
		if len(def) == 0 {
			return false, e
		}

		return def[0], nil
	}

	return b, e
}

func (c env) get(k string) (string, bool) {
	if c.useEnvs {
		cfg, exist := c.envs[k]
		return cfg, exist
	}

	cfg := os.Getenv(k)
	return cfg, cfg != ""
}

func (c env) GetStringArr(k string, def ...[]string) ([]string, error) {
	cfg, exist := c.get(k)
	if !exist {
		if len(def) == 0 {
			return nil, ErrorVarNotExist
		}

		return def[0], nil
	}

	return parser.String(cfg).ToStringArr()
}

func (c env) GetIntArr(k string, def ...[]int) ([]int, error) {
	cfg, exist := c.get(k)
	if !exist {
		if len(def) == 0 {
			return nil, ErrorVarNotExist
		}
		return def[0], nil
	}

	is, e := parser.String(cfg).ToIntArr()

	if e != nil {
		if len(def) == 0 {
			return nil, e
		}

		return def[0], nil
	}

	return is, e
}

func (c env) Dump() (string, error) {
	if !c.useEnvs {
		return "", errors.New("cannot dump env, you are using system-wide env")
	}

	j, e := json.Marshal(c.envs)

	return string(j), e
}

func LoadEnv(envPath string) (Env, error) {
	envs, err := godotenv.Read(envPath)

	if err == nil {
		for key, val := range envs {
			err = os.Setenv(key, val)
		}
	}

	return env{
		envs:    envs,
		useEnvs: err == nil,
	}, err
}
