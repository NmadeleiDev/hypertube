package dbInit

import (
	"auth_backend/logger"
	"auth_backend/postgres"
	"auth_backend/errors"
)

func RecreateTables() *errors.Error {

	print("Сбрасываю все таблицы\t\t\t- ")
	if Err := postgres.DropAllTables(); Err != nil {
		println(logger.RED + "ошибка" + logger.NO_COLOR)
		println(logger.RED + Err.Error() + logger.NO_COLOR)
		return Err
	}
	println(logger.GREEN + "успешно" + logger.NO_COLOR)

	print("Создаю таблицу базовых пользователей\t- ")
	if Err := postgres.CreateUsersTable(); Err != nil {
		println(logger.RED + "ошибка" + logger.NO_COLOR)
		println(logger.RED + Err.Error() + logger.NO_COLOR)
		return Err
	}
	println(logger.GREEN + "успешно" + logger.NO_COLOR)

	print("Создаю таблицу пользователей школы 42\t- ")
	if Err := postgres.CreateUsers42StrategyTable(); Err != nil {
		println(logger.RED + "ошибка" + logger.NO_COLOR)
		println(logger.RED + Err.Error() + logger.NO_COLOR)
		return Err
	}
	println(logger.GREEN + "успешно" + logger.NO_COLOR)

	print("Создаю таблицу пользователей vk.com\t- ")
	if Err := postgres.CreateUsersVkStrategyTable(); Err != nil {
		println(logger.RED + "ошибка" + logger.NO_COLOR)
		println(logger.RED + Err.Error() + logger.NO_COLOR)
		return Err
	}
	println(logger.GREEN + "успешно" + logger.NO_COLOR)

	print("Создаю таблицу пользователей facebook\t- ")
	if Err := postgres.CreateUsersFbStrategyTable(); Err != nil {
		println(logger.RED + "ошибка" + logger.NO_COLOR)
		println(logger.RED + Err.Error() + logger.NO_COLOR)
		return Err
	}
	println(logger.GREEN + "успешно" + logger.NO_COLOR)
	return nil
}
