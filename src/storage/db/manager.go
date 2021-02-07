package db

import (
	"hypertube_storage/dao"
	"hypertube_storage/db/postgres"
)

func GetLoadedFilesManager() dao.LoadedFilesDbManager {
	return &postgres.Manager
}