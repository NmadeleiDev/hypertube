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
	"hypertube_storage/subtitlesManager"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	videoRequest     = "video"
	subtitlesRequest = "srt"
)

func UploadFilePartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fileId := mux.Vars(r)["file_id"]
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
		logrus.Debugf("Got request with start=%v. Info: %v",
			fileRange.Start, info)

		if (info.IsLoaded || info.InProgress) && fileRange.Start >= info.VideoFile.Length {
			SendFailResponseWithCode(w, fmt.Sprintf(
				"Start byte (%v) in Content-Range exceeds file length (%v); info: %v",
					fileRange.Start, info.VideoFile.Length, info), http.StatusBadRequest)
			return
		}

		logrus.Debugf("Trying to upload video from %v", info)
		var filePart []byte

		if info.IsLoaded {
			filePart, _, err = filesReader.GetManager().GetFileInRange(info.VideoFile.Name, fileRange.Start, info.VideoFile.Length)
		} else if info.InProgress {
			filePart, _, err = filesReader.GetManager().GetFileInRange(info.VideoFile.Name, fileRange.Start, info.VideoFile.Length)
			if !filesReader.GetManager().IsPartWritten(info.VideoFile.Name, filePart, fileRange.Start) || err != nil {
				db.GetLoadedStateDb().PubPriorityByteIdx(fileId, info.VideoFile.Name, fileRange.Start)

				readCtx, readCancel := context.WithTimeout(context.TODO(), time.Second * 1800)
				defer readCancel()

				logrus.Debugf("Got file inProgress=true from db: %v, waiting for data (%v %v %v)", info.VideoFile.Name, filesReader.GetManager().HasNullBytes(filePart), filePart == nil, err)
				filePart, _, err = filesReader.GetManager().WaitForFilePart(readCtx, info.VideoFile.Name, fileRange.Start, info.VideoFile.Length)
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
			info, err = db.GetLoadedFilesManager().GetFileInfoById(fileId)
			if err != nil {
				SendFailResponseWithCode(w,fmt.Sprintf("File %s not found by id: %s", fileId, err.Error()), http.StatusNotFound)
				return
			}
			readCtx, readCancel := context.WithTimeout(context.TODO(), time.Second * 600)
			defer readCancel()

			db.GetLoadedStateDb().PubPriorityByteIdx(fileId, info.VideoFile.Name, fileRange.Start)

			filePart, _, err = filesReader.GetManager().WaitForFilePart(readCtx, info.VideoFile.Name, fileRange.Start, info.VideoFile.Length)
		}
		logrus.Debugf("Writing response, part len=%v", len(filePart))

		if err != nil {
			SendFailResponseWithCode(w, err.Error(), http.StatusInternalServerError)
		} else {
			fileRange.End = fileRange.Start + int64(len(filePart)) - 1
			contentLen := int64(len(filePart))
			w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", fileRange.Start, fileRange.End, info.VideoFile.Length))
			w.Header().Set("Content-Length", fmt.Sprint(contentLen))
			w.Header().Set("Accept-Ranges", "bytes")
			w.Header().Set("Content-Type", GetContentTypeForReqType("video"))
			w.WriteHeader(GetResponseStatusForReqType("video"))
			if _, err := io.Copy(w, bytes.NewReader(filePart)); err != nil {
				logrus.Errorf("Error piping response: %v", err)
			}

			go db.GetLoadedFilesManager().UpdateLastWatchedDate(fileId)
		}
	} else {
		SendFailResponseWithCode(w, "Incorrect method", http.StatusMethodNotAllowed)
	}
}

func UploadSubtitlesFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fileId := mux.Vars(r)["file_id"]
		subtitlesId := mux.Vars(r)["subtitles_id"]

		info, err := db.GetLoadedFilesManager().GetFileInfoById(fileId)
		if err != nil {
			logrus.Errorf("Err loading file '%v' info, err: %v", fileId, err)
			SendFailResponseWithCode(w,fmt.Sprintf("File %s not found by id: %s", fileId, err.Error()), http.StatusNotFound)
			return
		}

		var subtitlesFile []byte

		if info.IsLoaded {
			subtitlesFile, err = filesReader.GetManager().ReadWholeFile(subtitlesId)
		} else if info.InProgress {
			subtitlesFile, err = filesReader.GetManager().ReadWholeFile(subtitlesId)
			if !filesReader.GetManager().IsPartWritten(subtitlesId, subtitlesFile, 0) || err != nil {
				db.GetLoadedStateDb().PubPriorityByteIdx(fileId, subtitlesId, 0)

				readCtx, readCancel := context.WithTimeout(context.TODO(), time.Second * 600)
				defer readCancel()

				logrus.Debugf("Got file inProgress=true from db: %v, waiting for data (%v %v %v %v)", info, len(subtitlesFile), filesReader.GetManager().HasNullBytes(subtitlesFile), subtitlesFile == nil, err)
				subtitlesFile, err = filesReader.GetManager().WaitForWholeFileWritten(readCtx, subtitlesId)
				logrus.Debugf("Wait srt success (len=%v) from disk", len(subtitlesFile))
			} else {
				logrus.Debugf("Read srt file (len=%v) from disk", len(subtitlesFile))
			}
		} else {
			ok := SendTaskToTorrentClient(fileId)
			if !ok {
				SendFailResponseWithCode(w, "Failed to call torrent client", http.StatusInternalServerError)
				return
			}
			info, err = db.GetLoadedFilesManager().GetFileInfoById(fileId)
			if err != nil {
				SendFailResponseWithCode(w,fmt.Sprintf("File %s not found by id: %s", fileId, err.Error()), http.StatusNotFound)
				return
			}
			readCtx, readCancel := context.WithTimeout(context.TODO(), time.Second * 600)
			defer readCancel()

			db.GetLoadedStateDb().PubPriorityByteIdx(fileId, subtitlesId, 0)

			subtitlesFile, err = filesReader.GetManager().WaitForWholeFileWritten(readCtx, subtitlesId)
		}
		logrus.Debugf("Writing response, part len=%v", len(subtitlesFile))

		if err != nil {
			SendFailResponseWithCode(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.Header().Set("Content-Type", GetContentTypeForReqType("srt"))
			w.WriteHeader(GetResponseStatusForReqType("srt"))
			if err := subtitlesManager.GetManager().ConvertSrtToVtt(bytes.NewReader(subtitlesFile), w); err != nil {
				SendFailResponseWithCode(w, fmt.Sprintf("Failed to convert srt to vtt: %v", err.Error()), http.StatusInternalServerError)
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
