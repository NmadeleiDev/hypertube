package hash

import (
	"auth_backend/errors"
)

type Config struct {
	PasswdSalt string `conf:"passwordSalt"`
	MasterKey  string `conf:"masterKey"`
}

var cfg *Config

func GetConfig() *Config {
	if cfg == nil {
		cfg = &Config{}
	}
	return cfg
}

func getConfig() (*Config, *errors.Error) {
	if cfg == nil {
		return nil, errors.NotConfiguredPackage.SetArgs("controller/hash", "controller/hash")
	}
	return cfg, nil
}
