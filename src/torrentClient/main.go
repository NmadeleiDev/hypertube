package main

import (
	"torrentClient/db"
	"torrentClient/fsWriter"
	"torrentClient/parser/env"
	"torrentClient/server"
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

	go fsWriter.GetWriter().StartWaitingForData()
	server.Start()
}
