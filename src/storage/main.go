package main

import (
	"hypertube_storage/db"
	"hypertube_storage/parser/env"
	"hypertube_storage/server"
	"hypertube_storage/server/handlers"
)

func main() {
	InitLog()
	db.GetLoadedFilesManager().InitConnection(env.GetParser().GetPostgresDbDsn())
	db.GetLoadedFilesManager().InitTables()

	db.GetLoadedStateDb().InitConnection()

	defer func() {
		db.GetLoadedFilesManager().CloseConnection()
		db.GetLoadedStateDb().CloseConnection()
	}()

	go restartNotFinishedLoads()

	server.Start()
}

func restartNotFinishedLoads() {
	ids := db.GetLoadedFilesManager().GetInProgressFileIds()
	if ids == nil {
		return
	}

	for _, id := range ids {
		handlers.SendTaskToTorrentClient(id)
	}
}
