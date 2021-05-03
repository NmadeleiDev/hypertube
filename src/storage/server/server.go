package server

import (
	"net/http"

	"hypertube_storage/parser/env"
	"hypertube_storage/server/handlers"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var devMode bool

func Start() {
	devMode = env.GetParser().IsDevMode()

	router := mux.NewRouter()

	router.HandleFunc("/load/{file_id}", handlers.UploadFilePartHandler)
	router.HandleFunc("/load/{file_id}/{request}", handlers.UploadFilePartHandler)
	router.PathPrefix("/").HandlerFunc(handlers.CatchAllHandler)

	logrus.Info("Listening localhost:2222")
	if err := http.ListenAndServe(":2222", router); err != nil {
		logrus.Fatal("Server err: ", err)
	}
}
