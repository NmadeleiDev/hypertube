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

const (
	videoRequest = "video"
	srtRequest = "srt"
)

func UploadFilePartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fileId := mux.Vars(r)["file_id"]
		reqType := mux.Vars(r)["request"]
		if reqType == "" { // TODO потом удалить
			reqType = videoRequest
		}
		if reqType != videoRequest && reqType != srtRequest {
			SendFailResponseWithCode(w, fmt.Sprintf("Unknown req type: %v", reqType), http.StatusBadRequest)
			return
		}

		fileRange := model.FileRangeDescription{}
		if err := fileRange.ParseHeader(r.Header.Get("range")); err != nil {
			SendFailResponseWithCode(w, err.Error(), http.StatusBadRequest)
			return
		}

		info, err := db.GetLoadedFilesManager().GetFileInfoById(fileId)
		if err != nil {
			logrus.Errorf("Err loading file '%v' info, err: %v", fileId, err)
			SendFailResponseWithCode(w,fmt.Sprintf("File %s not found by id: %s", fileId, err.Error()), http.StatusNotFound)
			return
		}
		fileName, fileLength := GetFileInfoForReqType(info, reqType)

		logrus.Debugf("Got request with start = %v. Info: %v",
			fileRange.Start, info)

		if (info.IsLoaded || info.InProgress) && fileRange.Start >= fileLength {
			SendFailResponseWithCode(w, fmt.Sprintf(
				"Start byte (%v) in Content-Range exceeds file length (%v); info: %v",
					fileRange.Start, fileLength, info), http.StatusBadRequest)
			return
		}

		logrus.Debugf("Trying to upload type %v from %v", reqType, info)
		var filePart []byte

		if info.IsLoaded {
			filePart, _, err = filesReader.GetManager().GetFileInRange(fileName, fileRange.Start, fileLength)
		} else if info.InProgress {
			filePart, _, err = filesReader.GetManager().GetFileInRange(fileName, fileRange.Start, fileLength)
			if !filesReader.GetManager().IsPartWritten(fileName, filePart, fileRange.Start) || err != nil {
				db.GetLoadedStateDb().PubPriorityByteIdx(fileId, fileName, fileRange.Start)

				readCtx, readCancel := context.WithTimeout(context.TODO(), time.Second * 600)
				defer readCancel()

				logrus.Debugf("Got file inProgress=true from db: %v, waiting for data (%v %v %v)", fileName, filesReader.GetManager().HasNullBytes(filePart), filePart == nil, err)
				filePart, _, err = filesReader.GetManager().WaitForFilePart(readCtx, fileName, fileRange.Start, fileLength)
				logrus.Debugf("Wait success (start=%v, len=%v) from disk", fileRange.Start, len(filePart))
			} else {
				logrus.Debugf("Read part (start=%v, len=%v) from disk", fileRange.Start, len(filePart))
			}
		} else {
			ok := SendTaskToTorrentClient(fileId)
			if !ok {
				SendFailResponseWithCode(w, "Failed to call torrent client", http.StatusInternalServerError)
				return
			}
			fileName, fileLength, err = LoadFileInfoFromDbForReqType(fileId, reqType)
			if err != nil {
				SendFailResponseWithCode(w,fmt.Sprintf("File %s not found by id: %s", fileId, err.Error()), http.StatusNotFound)
				return
			}
			readCtx, readCancel := context.WithTimeout(context.TODO(), time.Second * 600)
			defer readCancel()

			db.GetLoadedStateDb().PubPriorityByteIdx(fileId, fileName, fileRange.Start)

			logrus.Debugf("Got file name from client: %v, waiting for data", fileName)
			filePart, _, err = filesReader.GetManager().WaitForFilePart(readCtx, fileName, fileRange.Start, fileLength)
		}
		logrus.Debugf("Writing response, part len=%v", len(filePart))

		if err != nil {
			SendFailResponseWithCode(w, err.Error(), http.StatusInternalServerError)
		} else {
			fileRange.End = fileRange.Start + int64(len(filePart)) - 1
			contentLen := int64(len(filePart))
			w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", fileRange.Start, fileRange.End, fileLength))
			w.Header().Set("Content-Length", fmt.Sprint(contentLen))
			w.Header().Set("Accept-Ranges", "bytes")
			w.Header().Set("Content-Type", GetContentTypeForReqType(reqType))
			w.WriteHeader(GetResponseStatusForReqType(reqType))
			if _, err := io.Copy(w, bytes.NewReader(filePart)); err != nil {
				logrus.Errorf("Error piping response: %v", err)
			}

			go db.GetLoadedFilesManager().UpdateLastWatchedDate(fileId)
		}
	} else {
		SendFailResponseWithCode(w, "Incorrect method", http.StatusMethodNotAllowed)
	}
}

func CatchAllHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("Catchall: %v", *r)
	SendFailResponseWithCode(w, "catchall", http.StatusNotFound)
}
