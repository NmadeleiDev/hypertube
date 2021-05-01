package main

import (
	"torrentClient/db"
	"torrentClient/fsWriter"
	"torrentClient/parser/env"
	"torrentClient/server"
	"torrentClient/torrentfile"

	"github.com/sirupsen/logrus"
)

func main() {
	InitLog()

	db.GetFilesManagerDb().InitConnection(env.GetParser().GetPostgresDbDsn())
	db.GetFilesManagerDb().InitTables()

	db.GetLoadedStateDb().InitConnection()

	defer func() {
		db.GetFilesManagerDb().CloseConnection()
		db.GetLoadedStateDb().CloseConnection()
	}()

	go restartInProgressLoads()
	go fsWriter.GetWriter().StartWaitingForData()

	server.Start()
}

func restartInProgressLoads()  {
	fileIds, ok := db.GetFilesManagerDb().GetInProgressFileIds()
	if !ok {
		return
	}

	for _, fileId := range fileIds {
		torrent, err := torrentfile.GetManager().LoadTorrentFileFromDB(fileId)
		if err != nil {
			logrus.Errorf("Error loading torrent: %v", err)
			return
		}

		torrent.PrepareDownload()

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
	}
}
