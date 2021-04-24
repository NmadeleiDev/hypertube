package db

import (
	"hypertube_storage/dao"
	"hypertube_storage/db/postgres"
	"hypertube_storage/db/redis"
)

func GetLoadedFilesManager() dao.LoadedFilesDbManager {
	return &postgres.Manager
}

func GetLoadedStateDb() dao.LoaderStateDbManager  {
	return &redis.Manager
}