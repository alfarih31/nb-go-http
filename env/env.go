package env

import (
	"encoding/json"
	"github.com/alfarih31/nb-go-http"
	"github.com/alfarih31/nb-go-http/parser"
	"github.com/joho/godotenv"
	"os"
)

type env struct {
	envs    map[string]string
	useEnvs bool
}

type Env interface {
	GetInt(k string, def int) (int, error)
	GetString(k string, def string) (string, error)
	GetBool(k string, def bool) (bool, error)
	GetStringArr(k string, def []string) ([]string, error)
	GetIntArr(k string, def []int) ([]int, error)
	Dump() (string, error)
}

func (c env) GetInt(k string, def int) (int, error) {
	cfg, exist := c.get(k)

	if !exist {
		return def, nil
	}

	i, e := parser.String(cfg).ToInt()

	if e != nil {
		return def, e
	}

	return i, e
}

func (c env) GetString(k string, def string) (string, error) {
	cfg, exist := c.get(k)
	if !exist {
		return def, nil
	}

	return cfg, nil
}

func (c env) GetBool(k string, def bool) (bool, error) {
	cfg, exist := c.get(k)
	if !exist {
		return def, nil
	}

	b, e := parser.String(cfg).ToBool()
	if e != nil {
		return def, e
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

func (c env) GetStringArr(k string, def []string) ([]string, error) {
	cfg, exist := c.get(k)
	if !exist {
		return def, nil
	}

	return parser.String(cfg).ToStringArr()
}

func (c env) GetIntArr(k string, def []int) ([]int, error) {
	cfg, exist := c.get(k)
	if !exist {
		return def, nil
	}

	is, e := parser.String(cfg).ToIntArr()

	if e != nil {
		return def, e
	}

	return is, e
}

func (c env) Dump() (string, error) {
	if !c.useEnvs {
		return "", noob.Err{Message: "Cannot dump env, you are using system-wide env!"}
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
