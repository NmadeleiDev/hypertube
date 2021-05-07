package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"torrentClient/db"
	"torrentClient/loadMaster"
	"torrentClient/magnetToTorrent"
	"torrentClient/torrentfile"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func DownloadRequestsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fileId := mux.Vars(r)["file_id"]

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
		if err != nil || torrent == nil {
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
			err = torrent.DownloadToFile()
			if err != nil {
				logrus.Errorf("Error downloading to file: %v", err)
			}
		}()
		SendDataResponse(w, response)
	}
}

func SubtitlesInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fileId := mux.Vars(r)["file_id"]

		type SubtitleRecord struct {
			Language string `json:"language"`
			Id       string `json:"id"`
			FileName	string	`json:"fileName"`
		}

		torrent, err := torrentfile.GetManager().LoadTorrentFileFromDB(fileId)
		if err != nil || torrent == nil {
			logrus.Errorf("Error loading torrent: %v", err)
			SendFailResponseWithCode(w, fmt.Sprintf("Error loading torrent from db: %v", err), http.StatusBadRequest)
			return
		}

		files := torrent.GetFiles()
		result := make([]SubtitleRecord, 0, len(files))
		for _, file := range files {
			if file.Extension() == "srt" {
				result = append(result, SubtitleRecord{
					Id: file.EncodeFileName(),
					Language: "unknown",
					FileName: strings.Join(file.Path, "_"),
				})
			}
		}

		SendDataResponse(w, result)
	}
}

func WriteLoadedPartsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fileId := mux.Vars(r)["file_id"]

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
		if err != nil || torrent == nil {
			logrus.Errorf("Error reading torrent file: %v", err)
			SendFailResponseWithCode(w, fmt.Sprintf("Error reading body: %s; body: %s", err.Error(), string(torrentBytes)), http.StatusInternalServerError)
			return
		}
		torrent.SysInfo.FileId = fileId

		logs := torrent.SaveLoadedPiecesToFS()
		SendDataResponse(w, logs)
	}
}

func LoadingStatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fileId := mux.Vars(r)["file_id"]

		stats, ok := loadMaster.GetMaster().GetStatsForEntry(fileId)
		logrus.Debugf("Got (ok=%v) %v stats: %v", ok, fileId, stats)
		if !ok {
			SendFailResponseWithCode(w, "Load not found", http.StatusBadRequest)
		} else {
			SendDataResponse(w, struct {
				ActivePeers	int `json:"activePeers"`
				LoadedPercent	int `json:"loadedPercent"`
				DonePieces	[]int `json:"donePieces"`
				InProgressPieces	[]int `json:"inProgressPieces"`
			}{
				ActivePeers: stats.NumOfActivePeers,
				LoadedPercent: len(stats.DonePieces) * 100 / stats.TotalPieces,
				DonePieces: stats.DonePieces,
				InProgressPieces: stats.InProgressPieces,
			})
		}
	} else {
		SendFailResponseWithCode(w, "Not allowed", http.StatusMethodNotAllowed)
	}
}

func TerminateLoadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fileId := mux.Vars(r)["file_id"]

		ok := loadMaster.GetMaster().StopLoad(fileId)
		if !ok {
			SendFailResponseWithCode(w, fmt.Sprintf("Load '%v' not found", fileId), http.StatusBadRequest)
		} else {
			SendSuccessResponse(w)
		}
	} else {
		SendFailResponseWithCode(w, "Not allowed", http.StatusMethodNotAllowed)
	}
}

