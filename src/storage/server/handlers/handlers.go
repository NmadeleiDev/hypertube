package handlers

import (
	"fmt"
	"io"
	"net/http"

	"hypertube_storage/db"
	"hypertube_storage/filesReader"
	"hypertube_storage/model"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)


func UploadFilePartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fileId := mux.Vars(r)["file_id"]
		fileRange := model.FileRangeDescription{}
		if err := fileRange.ParseHeader(r.Header.Get("range")); err != nil {
			SendFailResponseWithCode(w, err.Error(), http.StatusBadRequest)
			return
		}

		fileName, err := db.GetLoadedFilesManager().GetFilePathById(fileId)
		if err != nil {
			SendFailResponseWithCode(w,fmt.Sprintf("File %s not found: %s", fileId, err.Error()), http.StatusBadRequest)
			// вызов torrentClient для загрузки, ожидание загрузки нужного куска
		} else {
			filePart, totalLength, err := filesReader.GetManager().GetFileInRange(fileName, &fileRange)
			if err != nil {
				SendFailResponseWithCode(w, err.Error(), http.StatusInternalServerError)
			} else {
				w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", fileRange.Start, fileRange.End, totalLength))
				w.Header().Set("Content-Length", fmt.Sprint(totalLength))
				if _, err := io.Copy(w, filePart); err != nil {
					logrus.Errorf("Error piping response: %v", err)
				}
			}
		}
	}
}
