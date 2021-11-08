package nbgohttp

import (
	"github.com/joho/godotenv"
	"os"
)

type Env struct {
	envs    map[string]string
	useEnvs bool
}

func (c Env) GetInt(k string, def int) (int, error) {
	cfg, exist := c.get(k)

	if !exist {
		return def, nil
	}

	i, e := StringParser{cfg}.ToInt()

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

	b, e := StringParser{cfg}.ToBool()
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

	return StringParser{cfg}.ToStringArr()
}

func (c Env) GetIntArr(k string, def []int) ([]int, error) {
	cfg, exist := c.get(k)
	if !exist {
		return def, nil
	}

	is, e := StringParser{cfg}.ToIntArr()

	if e != nil {
		return def, e
	}

	return is, e
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
