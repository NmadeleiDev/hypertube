package dbInit

import (
	"auth_backend/logger"
	"auth_backend/postgres"
)

func InitDB() {
	//if Err := initer.InitPackages("conf.json"); Err != nil {
	//	println(logger.RED + Err.Error() + logger.NO_COLOR)
	//	return
	//}

	//defer func() {
	//	if Err := initer.CloseAllConnections(); Err != nil {
	//		println(logger.RED + "\nошибка при попытке закрытия соединения " + logger.NO_COLOR + Err.Error())
	//	} else {
	//		println(logger.GREEN + "\nвсе соединения с внешними службами успешно закрыты" + logger.NO_COLOR)
	//	}
	//}()

	//print("Сбрасываю все таблицы\t\t\t- ")
	//if Err := postgres.DropAllTables(); Err != nil {
	//	println(logger.RED + "ошибка" + logger.NO_COLOR)
	//	println(logger.RED + Err.Error() + logger.NO_COLOR)
	//	return
	//}
	//println(logger.GREEN + "успешно" + logger.NO_COLOR)

	print("Создаю таблицу базовых пользователей\t- ")
	if Err := postgres.CreateUsersTable(); Err != nil {
		println(logger.RED + "ошибка" + logger.NO_COLOR)
		println(logger.RED + Err.Error() + logger.NO_COLOR)
		return
	}
	println(logger.GREEN + "успешно" + logger.NO_COLOR)

	print("Создаю таблицу пользователей школы 42\t- ")
	if Err := postgres.CreateUsers42StrategyTable(); Err != nil {
		println(logger.RED + "ошибка" + logger.NO_COLOR)
		println(logger.RED + Err.Error() + logger.NO_COLOR)
		return
	}
	println(logger.GREEN + "успешно" + logger.NO_COLOR)

	print("Создаю таблицу пользователей vk.com\t- ")
	if Err := postgres.CreateUsersVkStrategyTable(); Err != nil {
		println(logger.RED + "ошибка" + logger.NO_COLOR)
		println(logger.RED + Err.Error() + logger.NO_COLOR)
		return
	}
	println(logger.GREEN + "успешно" + logger.NO_COLOR)

	print("Создаю таблицу пользователей facebook\t- ")
	if Err := postgres.CreateUsersFbStrategyTable(); Err != nil {
		println(logger.RED + "ошибка" + logger.NO_COLOR)
		println(logger.RED + Err.Error() + logger.NO_COLOR)
		return
	}
	println(logger.GREEN + "успешно" + logger.NO_COLOR)
}
