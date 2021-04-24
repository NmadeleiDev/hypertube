package db

import (
	"os"

	"torrentClient/db/postgres"
	"torrentClient/db/redis"
)

type LoaderStateDbManager interface {
	InitConnection()
	CloseConnection()

	AddSliceIndexForFile(fileName string, sliceByteIdx ...int64)
	DeleteSliceIndexesSet(fileName string)
}

func GetLoadedStateDb() LoaderStateDbManager  {
	return &redis.Manager
}

type FilesDbManager interface {
	InitTables()
	InitConnection(connStr string)
	CloseConnection()

	SetFileNameForRecord(fileId, name string)
	SetFileLengthForRecord(fileId string, length int64)
	SetInProgressStatusForRecord(fileId string, status bool)
	SetLoadedStatusForRecord(fileId string, status bool)

	GetLoadedIndexesForFile(fileId string) []int
	SaveFilePartsToFile(dest *os.File, fileId string, start int, length int) error
	LoadPartsForFile(fileId string, writeChan chan []byte)
	GetTorrentOrMagnetForByFileId(fileId string) ([]byte, string, bool)

	PreparePlaceForFile(fileId string)
	SaveFilePart(fileId string, part []byte, start, size, index int64)
	RemoveFilePartsPlace(fileId string)

	PartsTableNameForFile(fileId string) string
}

func GetFilesManagerDb() FilesDbManager {
	return &postgres.Manager
}