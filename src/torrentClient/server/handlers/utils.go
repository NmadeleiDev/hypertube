package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

	"torrentClient/model"

	"github.com/sirupsen/logrus"
)


func SendFailResponseWithCode(w http.ResponseWriter, text string, code int) {
	var packet []byte
	var err error

	response := &model.DataResponse{Status: false, Data: text}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)

	if packet, err = json.Marshal(response); err != nil {
		logrus.Error("Error marshalling response: ", err)
	}
	if _, err = w.Write(packet); err != nil {
		logrus.Error("Error sending response: ", err)
	}
}

func SendSuccessResponse(w http.ResponseWriter) {
	var packet []byte
	var err error

	response := &model.DataResponse{Status: true, Data: nil}
	w.Header().Set("content-type", "application/json")

	if packet, err = json.Marshal(response); err != nil {
		logrus.Error("Error marshalling response: ", err)
	}
	if _, err = w.Write(packet); err != nil {
		logrus.Error("Error sending response: ", err)
	}
}

func SendDataResponse(w http.ResponseWriter, data interface{}) {
	var packet []byte
	var err error

	response := &model.DataResponse{Status: true, Data: data}
	w.Header().Set("content-type", "application/json")

	if packet, err = json.Marshal(response); err != nil {
		logrus.Error("Error marshalling response: ", err)
	}
	if _, err = w.Write(packet); err != nil {
		logrus.Error("Error sending response: ", err)
	}
}

func SetCookieForHour(w http.ResponseWriter, cookieName, value string) {
	c := http.Cookie{
		Name:     cookieName,
		Domain: ".gomobile.ru",
		Value:    value,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(time.Hour.Seconds())}
	http.SetCookie(w, &c)
}

func GetTrackersFromMagnet(magnet string) string {
	decoded, err := url.ParseQuery(magnet)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return decoded.Get("tr")
}
