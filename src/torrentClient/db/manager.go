package db

import (
	"os"

	"torrent_client/db/postgres"
	"torrent_client/db/redis"
)

type LoaderStateDbManager interface {
	InitConnection()
	CloseConnection()

	GetActiveDownloads() []string
	CheckIfFileIsActiveLoading(file string) bool
	AddFileIdToActiveDownloads(id string)
	AnnounceLoadedPart(fileId, partId string, start, size int64)
	SaveLoadedPartInfo(fileId, partId string, start, size int64)

	CleanLoadingLogsForFile(file string)
	DeleteFileFromActiveDownloads(file string)
}

func GetLoadedStateDb() LoaderStateDbManager  {
	return &redis.Manager
}

type FilesDbManager interface {
	InitTables()
	InitConnection(connStr string)
	CloseConnection()

	SaveFileNameForReadyFile(fileId, name string)
	SaveFilePartsToFile(dest *os.File, fileId string)
	GetTorrentFileForByFileId(fileId string) []byte

	PreparePlaceForFile(fileId string)
	SaveFilePart(fileId string, part []byte, start, size, index int64)
	RemoveFilePartsPlace(fileId string)
}

func GetFilesManagerDb() FilesDbManager {
	return &postgres.Manager
}