package controller

import (
	"auth_backend/errors"
)

type Config struct {
	ProjectRoot          string `conf:"projectRoot"`
	ServerPasswd         string `conf:"passwd"`
	ServerIp             string `conf:"ip"`
	ServerPort           uint   `conf:"port"`
	SocketRedirect       string `conf:"socketRedirect"`
	OauthRedirect        string `conf:"oauthRedirect"`
	ErrorRedirect        string `conf:"errorRedirect"`
	PasswdResetRedirect  string `conf:"passwdResetRedirect"`
	Ecole42ClientId      string `conf:"ecole42ClientId"`
	Ecole42Secret        string `conf:"ecole42Secret"`
	FacebookClientId     string `conf:"facebookClientId"`
	FacebookSecret       string `conf:"facebookSecret"`
	VkClientId           string `conf:"vkClientId"`
	VkSecret             string `conf:"vkSecret"`
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
		return nil, errors.NotConfiguredPackage.SetArgs("controller", "controller")
	}
	return cfg, nil
}

func GetServerPort() (uint, *errors.Error) {
	conf, Err := getConfig()
	if Err != nil {
		return 0, Err
	}
	return conf.ServerPort, nil
}
