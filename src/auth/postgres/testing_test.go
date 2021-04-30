package postgres

import (
	"auth_backend/configurator"
	"auth_backend/errors"
	"auth_backend/logger"
	"testing"
)

var (
	emailValid1   = "user1@gmail.com"
	emailValid2   = "user2@gmail.com"
	encryptedPass = "WQEsafqwesa="
	username      = "USER"
)

func initTest(t *testing.T) {
	if Err := initPackages("../conf.json"); Err != nil {
		t.Errorf("%sError: cannot init packages - %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		t.FailNow()
	}
}

func initPackages(configFileName string) *errors.Error {
	/*
	**	Read config
	 */
	print("Считываю конфигурационный файл\t\t- ")
	if err := configurator.SetConfigFile(configFileName); err != nil {
		println(logger.RED + "ошибка" + logger.NO_COLOR)
		return errors.ConfigurationFail.SetArgs("Не могу считать файл "+configFileName,
			"Cant read file "+configFileName).SetOrigin(err)
	}
	println(logger.GREEN + "успешно" + logger.NO_COLOR)
	/*
	**	logger
	 */
	print("Настраиваю пакет logger\t\t\t- ")
	cfgLogger := logger.GetConfig()
	if err := configurator.ParsePackageConfig(cfgLogger, "logger"); err != nil {
		println(logger.RED + "ошибка" + logger.NO_COLOR)
		return errors.ConfigurationFail.SetArgs("logger", "logger").SetOrigin(err)
	}
	println(logger.GREEN + "успешно" + logger.NO_COLOR)
	/*
	**	postrgres
	 */
	print("Настраиваю пакет postgres\t\t- ")
	cfgPostgres := GetConfig()
	if err := configurator.ParsePackageConfig(cfgPostgres, "postgres"); err != nil {
		println(logger.RED + "ошибка" + logger.NO_COLOR)
		return errors.ConfigurationFail.SetArgs("postgres", "postgres").SetOrigin(err)
	}
	println(logger.GREEN + "успешно" + logger.NO_COLOR)
	/*
	**	connection create
	 */
	print("Инициализирую соединение с postgres\t- ")
	var Err *errors.Error
	if Err = Init(); Err != nil {
		println(logger.RED + "ошибка" + logger.NO_COLOR)
		return Err
	}
	println(logger.GREEN + "успешно" + logger.NO_COLOR)
	return nil
}
