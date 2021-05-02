package postgres

import (
	"auth_backend/errors"
)

type Config struct {
	Host           string `conf:"host"`
	User           string `conf:"user"`
	Passwd         string `conf:"passwd"`
	Database       string `conf:"databaseName"`
	Type           string `conf:"databaseType"`
	ConnMax        uint   `conf:"connectionsMax"`
	RecreateTables bool   `conf:"recreateTables"`
}

var cfg *Config

func GetConfig() *Config {
	if cfg == nil {
		cfg = &Config{}
	}
	return cfg
}

func getConnection() (*Connection, *errors.Error) {
	if cfg == nil {
		return nil, errors.NotConfiguredPackage.SetArgs("postgres", "postgres")
	}
	if gConnection == nil {
		return nil, errors.NotInitializedPackage.SetArgs("postgres", "postgres")
	}
	return gConnection, nil
}
