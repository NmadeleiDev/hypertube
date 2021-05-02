package eraser

import (
	"time"

	"hypertube_storage/dao"
	"hypertube_storage/db"
	"hypertube_storage/filesReader"
)

type EraseManager struct {
}

var manager EraseManager

func GetEraser() dao.RecordsEraseManager {
	return &manager
}

func (e *EraseManager) StartCheckingForRecords()  {
	ticker := time.NewTicker(time.Hour * 3)

	for {
		<- ticker.C

		unwatchedIds, names := db.GetLoadedFilesManager().GetFileIdsWithWatchedUnder(time.Now().AddDate(0, -1, 0))

		for i, id := range unwatchedIds {
			if err := db.GetLoadedFilesManager().DeleteLoadedFileInfo(id); err != nil {
				continue
			}

			filesReader.GetManager().RemoveFile(names[i])
		}
	}
}
