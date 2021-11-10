package env

import (
	"encoding/json"
	"github.com/alfarih31/nb-go-http/app_error"
	"github.com/alfarih31/nb-go-http/data"
	"github.com/joho/godotenv"
	"os"
)

type Env struct {
	envs    map[string]string
	useEnvs bool
}

type IEnv interface {
	GetInt(k string, def int) (int, error)
	GetString(k string, def string) (string, error)
	GetBool(k string, def bool) (bool, error)
	GetStringArr(k string, def []string) ([]string, error)
	GetIntArr(k string, def []int) ([]int, error)
	Dump() (string, error)
}

func (c Env) GetInt(k string, def int) (int, error) {
	cfg, exist := c.get(k)

	if !exist {
		return def, nil
	}

	i, e := data.StringParser{cfg}.ToInt()

	if e != nil {
		return def, e
	}

	return i, e
}

func (c Env) GetString(k string, def string) (string, error) {
	cfg, exist := c.get(k)
	if !exist {
		return def, nil
	}

	return cfg, nil
}

func (c Env) GetBool(k string, def bool) (bool, error) {
	cfg, exist := c.get(k)
	if !exist {
		return def, nil
	}

	b, e := data.StringParser{cfg}.ToBool()
	if e != nil {
		return def, e
	}

	return b, e
}

func (c Env) get(k string) (string, bool) {
	if c.useEnvs {
		cfg, exist := c.envs[k]
		return cfg, exist
	}

	cfg := os.Getenv(k)
	return cfg, cfg != ""
}

func (c Env) GetStringArr(k string, def []string) ([]string, error) {
	cfg, exist := c.get(k)
	if !exist {
		return def, nil
	}

	return data.StringParser{cfg}.ToStringArr()
}

func (c Env) GetIntArr(k string, def []int) ([]int, error) {
	cfg, exist := c.get(k)
	if !exist {
		return def, nil
	}

	is, e := data.StringParser{cfg}.ToIntArr()

	if e != nil {
		return def, e
	}

	return is, e
}

func (c Env) Dump() (string, error) {
	if !c.useEnvs {
		return "", apperror.Err{Message: "Cannot dump env, you are using system-wide env!"}
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

	return Env{
		envs:    envs,
		useEnvs: err == nil,
	}, err
}
