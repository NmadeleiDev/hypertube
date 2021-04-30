package torrentfile

import (
	"context"

	"torrentClient/db"

	"github.com/sirupsen/logrus"
)

type LoadPriority struct {
	torrentFile	*TorrentFile
}

func (p *LoadPriority) StartPriorityUpdating(ctx context.Context) chan int64 {
	updatesChan := db.GetLoadedStateDb().GetLoadPriorityUpdatesChan(ctx, p.torrentFile.SysInfo.FileId)

	resultsChan := make(chan int64, 100)
	go func() {
		for {
			select {
			case <- ctx.Done():
				return
			case update := <- updatesChan:
				for _, file := range p.torrentFile.FileBoundariesMapping {
					if file.FileName == update.FileName {
						if update.ByteIdx > (file.End - file.Start) {
							logrus.Errorf("Wow! update.ByteIdx (%v) > (file.End - file.Start) (%v)", update.ByteIdx, file.End - file.Start)
							continue
						}
						overAllByteIdx := file.Start + update.ByteIdx
						pieceIdx := overAllByteIdx / int64(p.torrentFile.PieceLength)
						if pieceIdx > int64(len(p.torrentFile.PieceHashes)) {
							logrus.Errorf("Wow! pieceIdx (%v) > int64(len(p.torrentFile.PieceHashes)) (%v)", pieceIdx, len(p.torrentFile.PieceHashes))
							continue
						}
						resultsChan <- pieceIdx
					}
				}
			}
		}
	}()

	return resultsChan
}
