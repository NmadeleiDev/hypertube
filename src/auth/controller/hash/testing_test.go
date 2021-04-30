package hash

import (
	"auth_backend/configurator"
	"auth_backend/logger"
	"testing"
)

func initializePackageForTest(t *testing.T) {
	var configFileName = "../../conf.json"
	print("Считываю конфигурационный файл\t\t- ")
	if err := configurator.SetConfigFile(configFileName); err != nil {
		println(logger.RED + "ошибка" + logger.NO_COLOR)
		t.Errorf("%sНе могу считать файл %s: %s%s", logger.RED_BG, configFileName, err.Error(), logger.NO_COLOR)
		t.FailNow()
	}
	println(logger.GREEN + "успешно" + logger.NO_COLOR)
	/*
	**	logger
	 */
	print("Настраиваю пакет logger\t\t\t- ")
	cfgLogger := logger.GetConfig()
	if err := configurator.ParsePackageConfig(cfgLogger, "logger"); err != nil {
		println(logger.RED + "ошибка" + logger.NO_COLOR)
		t.Errorf("%sНе могу инициализировать пакет: %s%s", logger.RED_BG, err.Error(), logger.NO_COLOR)
		t.FailNow()
	}
	println(logger.GREEN + "успешно" + logger.NO_COLOR)
	/*
	**	hash
	 */
	print("Настраиваю пакет hash\t\t\t- ")
	cfgHash := GetConfig()
	if err := configurator.ParsePackageConfig(cfgHash, "hash"); err != nil {
		println(logger.RED + "ошибка" + logger.NO_COLOR)
		t.Errorf("%sНе могу инициализировать пакет: %s%s", logger.RED_BG, err.Error(), logger.NO_COLOR)
		t.FailNow()
	}
	println(logger.GREEN + "успешно" + logger.NO_COLOR)
}
