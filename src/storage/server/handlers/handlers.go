package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"hypertube_storage/db"
	"hypertube_storage/filesReader"
	"hypertube_storage/model"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func UploadFilePartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fileId := mux.Vars(r)["file_id"]
		logrus.Debugf("Got req for id: %v", fileId)
		fileRange := model.FileRangeDescription{}
		if err := fileRange.ParseHeader(r.Header.Get("range")); err != nil {
			SendFailResponseWithCode(w, err.Error(), http.StatusBadRequest)
			return
		}

		fileName, inProgress, isLoaded, fileLength, err := db.GetLoadedFilesManager().GetFileInfoById(fileId)
		if err != nil {
			logrus.Errorf("File not found by id '%v', err: %v", fileId, err)
			SendFailResponseWithCode(w,fmt.Sprintf("File %s not found by id: %s", fileId, err.Error()), http.StatusNotFound)
			return
		}

		if fileRange.Start >= fileLength && isLoaded {
			SendFailResponseWithCode(w,
				fmt.Sprintf("Start byte (%v) in Content-Range exceeds file length (%v)",
					fileRange.Start, fileLength), http.StatusBadRequest)
			return
		}

		logrus.Debugf("Trying to upload file %v", fileName)

		var filePart []byte

		if isLoaded {
			filePart, _, err = filesReader.GetManager().GetFileInRange(fileName, fileRange.Start)
		} else if inProgress {
			filePart, _, err = filesReader.GetManager().GetFileInRange(fileName, fileRange.Start)
			if !filesReader.GetManager().IsPartWritten(fileName, filePart, fileRange.Start) || err != nil {
				db.GetLoadedStateDb().PubPriorityByteIdx(fileId, fileName, fileRange.Start)

				readCtx, readCancel := context.WithTimeout(context.TODO(), time.Second * 600)
				defer readCancel()

				logrus.Debugf("Got file inProgress=true from db: %v, waiting for data (%v %v %v)", fileName, filesReader.GetManager().HasNullBytes(filePart), filePart == nil, err)
				filePart, _, err = filesReader.GetManager().WaitForFilePart(readCtx, fileName, fileRange.Start)
			}
		} else {
			fileName, ok := SendTaskToTorrentClient(fileId)
			if !ok {
				SendFailResponseWithCode(w, "Failed to call torrent client", http.StatusInternalServerError)
				return
			}
			readCtx, readCancel := context.WithTimeout(context.TODO(), time.Second * 60)
			defer readCancel()

			db.GetLoadedStateDb().PubPriorityByteIdx(fileId, fileName, fileRange.Start)

			logrus.Debugf("Got file name from client: %v, waiting for data", fileName)
			filePart, _, err = filesReader.GetManager().WaitForFilePart(readCtx, fileName, fileRange.Start)
		}
		logrus.Debugf("Writing response: %v", len(filePart))

		if err != nil {
			SendFailResponseWithCode(w, err.Error(), http.StatusInternalServerError)
		} else {
			fileRange.End = fileRange.Start + int64(len(filePart))
			contentLen := fileRange.End - fileRange.Start
			w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", fileRange.Start, fileRange.End, fileLength))
			w.Header().Set("Content-Length", fmt.Sprint(contentLen))
			w.Header().Set("Accept-Ranges", "bytes")
			w.Header().Set("Content-Type", "video/mp4")
			w.WriteHeader(http.StatusPartialContent)
			if _, err := io.Copy(w, bytes.NewReader(filePart)); err != nil {
				logrus.Errorf("Error piping response: %v", err)
			}
		}
	} else {
		SendFailResponseWithCode(w, "Incorrect method", http.StatusMethodNotAllowed)
	}
}

func CatchAllHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("Catchall: %v", *r)
}
