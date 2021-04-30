package main

import (
	"auth_backend/controller"
	"auth_backend/initer"
	"auth_backend/logger"
	"auth_backend/dbInit"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	if Err := initer.InitPackages("conf.json"); Err != nil {
		println(logger.RED + Err.Error() + logger.NO_COLOR)
		return
	}

	dbInit.InitDB()

	mux := controller.Router()
	println("Настраиваю роутер\t\t\t- " + logger.GREEN + "успешно" + logger.NO_COLOR)

	portUint, Err := controller.GetServerPort()
	if Err != nil {
		println(logger.RED + Err.Error() + logger.NO_COLOR)
		return
	}
	portString := strconv.FormatUint(uint64(portUint), 10)

	go func(portString string) {
		println("стартую сервер аутентификации на :" + portString + "\t- " + logger.GREEN + "успешно" + logger.NO_COLOR)
		http.ListenAndServe(":"+portString, mux)
		println(logger.RED + "порт :" + portString + " занят" + logger.NO_COLOR)
	}(portString)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if Err := initer.CloseAllConnections(); Err != nil {
		println(logger.RED + "\nошибка при попытке закрытия соединения " + logger.NO_COLOR + Err.Error())
	} else {
		println(logger.GREEN + "\nвсе соединения с внешними службами успешно закрыты" + logger.NO_COLOR)
	}
}
