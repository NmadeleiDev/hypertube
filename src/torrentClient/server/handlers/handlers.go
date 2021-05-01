package handlers

import (
	"bytes"
	"fmt"
	"net/http"

	"torrentClient/db"
	"torrentClient/magnetToTorrent"
	"torrentClient/torrentfile"

	"github.com/sirupsen/logrus"
)

func DownloadRequestsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fileId := r.URL.Query().Get("file_id")

		response := struct {
			IsLoaded	bool	`json:"isLoaded"`
			IsLoading	bool	`json:"isLoading"`
			FileName	string		`json:"fileName"`
			FileLength	int			`json:"fileLength"`
		}{}

		inProgress, isLoaded, ok := db.GetFilesManagerDb().GetFileStatus(fileId)
		if !ok {
			SendFailResponseWithCode(w, "Not found torrent in db", http.StatusBadRequest)
		}

		torrent, err := torrentfile.GetManager().LoadTorrentFileFromDB(fileId)
		if err != nil {
			logrus.Errorf("Error loading torrent: %v", err)
			SendFailResponseWithCode(w, fmt.Sprintf("Error loading torrent from db: %v", err), http.StatusBadRequest)
			return
		}
		response.FileLength = torrent.GetVideoFileLength()

		if isLoaded || inProgress {
			response.IsLoaded = isLoaded
			response.IsLoading = inProgress
			response.FileName = torrent.GetVideoFileName()

			SendDataResponse(w, response)
			return
		}

		response.FileName, _ = torrent.PrepareDownload()

		go func() {
			db.GetFilesManagerDb().SetInProgressStatusForRecord(torrent.SysInfo.FileId, true)
			defer db.GetFilesManagerDb().SetInProgressStatusForRecord(torrent.SysInfo.FileId, false)

			err = torrent.DownloadToFile()
			if err != nil {
				logrus.Errorf("Error downloading to file: %v", err)
			} else {
				db.GetFilesManagerDb().SetLoadedStatusForRecord(torrent.SysInfo.FileId, true)
			}
		}()
		SendDataResponse(w, response)
	}
}

func WriteLoadedPartsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fileId := r.URL.Query().Get("file_id")

		torrentBytes, magnetLink, ok := db.GetFilesManagerDb().GetTorrentOrMagnetForByFileId(fileId)
		if !ok {
			SendFailResponseWithCode(w, "File not found or not downloadable", http.StatusNotFound)
			return
		}

		if (torrentBytes == nil || len(torrentBytes) == 0) && len(magnetLink) > 0 {
			torrentBytes = magnetToTorrent.ConvertMagnetToTorrent(magnetLink)
			logrus.Info("Converted! ", len(torrentBytes))
		}

		torrent, err := torrentfile.GetManager().ReadTorrentFileFromBytes(bytes.NewBuffer(torrentBytes))
		if err != nil {
			logrus.Errorf("Error reading torrent file: %v", err)
			SendFailResponseWithCode(w, fmt.Sprintf("Error reading body: %s; body: %s", err.Error(), string(torrentBytes)), http.StatusInternalServerError)
			return
		}
		torrent.SysInfo.FileId = fileId

		logs := torrent.SaveLoadedPiecesToFS()
		SendDataResponse(w, logs)
	}
}
