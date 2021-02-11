package main

import (
	"torrent_client/db"
	"torrent_client/parser/env"
	"torrent_client/server"
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

	server.Start()
}
