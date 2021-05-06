package dao

import (
	"context"
	"time"

	"hypertube_storage/model"
)

type LoadedFilesDbManager interface {
	InitConnection(connStr string)
	InitTables()
	CloseConnection()

	GetFileInfoById(id string) (info model.LoadInfo, err error)
	GetInProgressFileIds() (result []string)
	GetFileIdsWithWatchedUnder(under time.Time) (ids []string, names []string)
	UpdateLastWatchedDate(fileId string)
	DeleteLoadedFileInfo(id string) error
}

type FileReader interface {
	WaitForFilePart(ctx context.Context, fileName string, start int64, expectedLen int64) ([]byte, int64, error)
	WaitForWholeFileWritten(ctx context.Context, fileName string) ([]byte, error)
	GetFileInRange(fileName string, start int64, expectedLen int64) (result []byte, totalLength int64, err error)
	ReadWholeFile(fileName string) ([]byte, error)
	HasNullBytes(src []byte) bool
	HasNotNullBytes(src []byte) bool
	IsPartWritten(fileName string, part []byte, start int64) bool
	RemoveFile(fileName string) bool
}

type RecordsEraseManager interface {
	StartCheckingForRecords()
}

type LoaderStateDbManager interface {
	InitConnection()
	CloseConnection()

	GetSliceIndexesForFile(fileName string) []int64
	PubPriorityByteIdx(fileId, fileName string, idx int64)
}
