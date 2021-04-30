package dao

import (
	"context"
)

type LoadedFilesDbManager interface {
	InitConnection(connStr string)
	InitTables()
	CloseConnection()

	GetFileInfoById(id string) (path string, inProgress, isLoaded bool, fLen int64, err error)
	GetInProgressFileIds() (result []string)
}

type FileReader interface {
	WaitForFilePart(ctx context.Context, fileName string, start int64) ([]byte, int64, error)
	GetFileInRange(fileName string, start int64) (result []byte, totalLength int64, err error)
	HasNullBytes(src []byte) bool
	HasNotNullBytes(src []byte) bool
	IsPartWritten(fileName string, part []byte, start int64) bool
}

type LoaderStateDbManager interface {
	InitConnection()
	CloseConnection()

	GetSliceIndexesForFile(fileName string) []int64
	PubPriorityByteIdx(fileId, fileName string, idx int64)
}
