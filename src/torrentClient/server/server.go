package server

import (
	"net/http"

	"torrentClient/parser/env"
	"torrentClient/server/handlers"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var devMode bool

type UserData struct {
	UserEmail string `json:"Email"`
	UserToken string `json:"Token"`
}

type TokenRequest struct {
	Token string `json:"Token"`
}

func Start() {
	devMode = env.GetParser().IsDevMode()

	router := mux.NewRouter()

	router.HandleFunc("/download/{file_id}", handlers.DownloadRequestsHandler)
	router.HandleFunc("/save/{file_id}", handlers.WriteLoadedPartsHandler)

	router.HandleFunc("/stats/{file_id}", handlers.LoadingStatsHandler)
	router.HandleFunc("/subtitles/{file_id}", handlers.SubtitlesInfoHandler)
	router.HandleFunc("/stop/{file_id}", handlers.TerminateLoadHandler)

	logrus.Info("Listening localhost:2222")
	if err := http.ListenAndServe(":2222", router); err != nil {
		logrus.Fatal("Server err: ", err)
	}
}

