package eraser

import (
	"time"

	"torrentClient/db"
	"torrentClient/fsWriter"
	"torrentClient/torrentfile"

	"github.com/sirupsen/logrus"
)

type EraseManager struct {
}

var manager EraseManager

type RecordsEraseManager interface {
	StartCheckingForRecords()
}

func GetEraser() RecordsEraseManager {
	return &manager
}

func (e *EraseManager) StartCheckingForRecords()  {
	ticker := time.NewTicker(time.Second * 30)

	for {
		<- ticker.C

		unwatchedIds := db.GetFilesManagerDb().GetFileIdsWithWatchedUnder(time.Now().AddDate(0, -1, 0))

		for _, id := range unwatchedIds {
			logrus.Debugf("Deleting files for %v", id)
			torrent, err := torrentfile.GetManager().LoadTorrentFileFromDB(id)

			if err := db.GetFilesManagerDb().DeleteLoadedFileInfo(id); err != nil {
				logrus.Errorf("Error deleting loaded file info %v", err)
			}

			if err != nil || torrent == nil {
				logrus.Errorf("Error loading torrent: %v", err)
				continue
			}

			db.GetFilesManagerDb().RemoveFilePartsPlace(id)

			files := torrent.GetFiles()
			for _, file := range files {
				logrus.Debugf("Removing file %v", file.EncodeFileName())
				fsWriter.GetWriter().RemoveFile(file.EncodeFileName())
			}
		}
	}
}
