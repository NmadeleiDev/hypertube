package main

import (
	"hypertube_storage/db"
	"hypertube_storage/parser/env"
	"hypertube_storage/server"
)

func main() {
	InitLog()
	db.GetLoadedFilesManager().InitConnection(env.GetParser().GetPostgresDbDsn())
	db.GetLoadedFilesManager().InitTables()

	defer func() {
		db.GetLoadedFilesManager().CloseConnection()
	}()

	server.Start()
}
