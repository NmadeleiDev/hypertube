package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"hypertube_storage/db"
	"hypertube_storage/model"
	"hypertube_storage/parser/env"

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

func GetFileInfoForReqType(info model.LoadInfo, reqType string) (fileName string, fileLength int64) {
	if reqType == videoRequest {
		fileName = info.VideoFile.Name
		fileLength = info.VideoFile.Length
	} else {
		fileName = info.SrtFile.Name
		fileLength = info.SrtFile.Length
	}

	return fileName, fileLength
}

func LoadFileInfoFromDbForReqType(fileId string, reqType string) (fileName string, fileLength int64, err error) {
	info, err := db.GetLoadedFilesManager().GetFileInfoById(fileId)
	if err != nil {
		logrus.Errorf("File not found by id '%v', err: %v", fileId, err)
		return "", 0, err
	}
	fileName, fileLength = GetFileInfoForReqType(info, reqType)
	return fileName, fileLength, nil
}

func GetContentTypeForReqType(reqType string) string {
	switch reqType {
	case videoRequest:
		return "video/mp4"
	case srtRequest:
		return "application/octet-stream"
	default:
		return "plain/text"
	}
}

func GetResponseStatusForReqType(reqType string) int {
	switch reqType {
	case videoRequest:
		return http.StatusPartialContent
	case srtRequest:
		return http.StatusOK
	default:
		return http.StatusOK
	}
}

func SendTaskToTorrentClient(fileId string) bool {
	req, err := http.Get(fmt.Sprintf("http://%s/download/%s", env.GetParser().GetLoaderServiceHost(), fileId))
	if err != nil {
		logrus.Errorf("Error calling loader service: %v", err)
		return false
	}

	if req.StatusCode != http.StatusOK {
		logrus.Errorf("Not ok status from torrent client: %v %v", req.StatusCode, req.Status)
		return false
	}

	info := model.LoaderTaskResponse{}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		logrus.Errorf("Error reading body: %v", err)
		return false
	}

	if err := json.Unmarshal(body, &info); err != nil {
		logrus.Errorf("Error unmarshal body from loader: %v", err)
		return false
	}

	return true
}


