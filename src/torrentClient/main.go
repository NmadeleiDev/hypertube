package main

import (
	"torrentClient/db"
	"torrentClient/eraser"
	"torrentClient/fsWriter"
	"torrentClient/loadMaster"
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

	if env.GetParser().DoRestartInProgressLoads() {
		go restartInProgressLoads()
	}

	loadMaster.GetMaster().Init()

	go fsWriter.GetWriter().StartWaitingForData()

	go eraser.GetEraser().StartCheckingForRecords()

	server.Start()
}

func restartInProgressLoads()  {
	fileIds, ok := db.GetFilesManagerDb().GetInProgressFileIds()
	if !ok {
		return
	}

	for _, fileId := range fileIds {
		torrent, err := torrentfile.GetManager().LoadTorrentFileFromDB(fileId)
		if err != nil || torrent == nil {
			logrus.Errorf("Error loading torrent: %v", err)
			return
		}

		torrent.PrepareDownload()

		go func() {
			db.GetFilesManagerDb().SetInProgressStatusForRecord(torrent.SysInfo.FileId, true)

			err = torrent.DownloadToFile()
			if err != nil {
				logrus.Errorf("Error downloading to file: %v", err)
			} else {
				db.GetFilesManagerDb().SetInProgressStatusForRecord(torrent.SysInfo.FileId, false)
				db.GetFilesManagerDb().SetLoadedStatusForRecord(torrent.SysInfo.FileId, true)
			}
		}()
	}
}
