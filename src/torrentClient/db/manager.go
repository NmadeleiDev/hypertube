package db

import (
	"context"
	"os"

	"torrentClient/db/postgres"
	"torrentClient/db/redis"
)

type LoaderStateDbManager interface {
	InitConnection()
	CloseConnection()

	AddSliceIndexForFile(fileName string, sliceByteIdx ...int64)
	DeleteSliceIndexesSet(fileName string)
	GetLoadPriorityUpdatesChan(ctx context.Context, fileId string) chan redis.PriorityUpdateMsg
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
	SetVideoFileNameAndLengthForRecord(fileId, fileName string, length int64)
	SetSrtFileNameAndLengthForRecord(fileId, fileName string, length int64)
	SetInProgressStatusForRecord(fileId string, status bool)
	SetLoadedStatusForRecord(fileId string, status bool)
	GetFileStatus(fileId string) (inProgress bool, isLoaded bool, ok bool)
	GetInProgressFileIds() (fileIds []string, ok bool)

	GetLoadedIndexesForFile(fileId string) []int
	GetPartDataByIdx(fileId string, idx int) ([]byte, int64, int64, bool)
	SaveFilePartsToFile(dest *os.File, fileId string, start int, length int) error
	DropDataPartByIdx(fileId string, idx int) bool
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