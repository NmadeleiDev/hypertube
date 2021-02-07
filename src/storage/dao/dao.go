package dao

import (
	"io"

	"hypertube_storage/model"
)

type LoadedFilesDbManager interface {
	InitConnection(connStr string)
	InitTables()
	CloseConnection()

	GetFilePathById(id string) (path string, err error)
}

type FileReader interface {
	GetFileInRange(path string, description *model.FileRangeDescription) (reader io.Reader, totalLength int64, err error)
}
