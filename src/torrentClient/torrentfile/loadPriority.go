package torrentfile

import (
	"context"

	"torrentClient/db"

	"github.com/sirupsen/logrus"
)

type LoadPriority struct {
	torrentFile	*TorrentFile
}

func (p *LoadPriority) StartPriorityUpdating(ctx context.Context) chan int {
	inputChan := db.GetLoadedStateDb().GetLoadPriorityUpdatesChan(ctx, p.torrentFile.SysInfo.FileId)

	outputChan := make(chan int, 100)

	go func() {
		defer close(outputChan)

		nPieces := len(p.torrentFile.PieceHashes)
		pLen := p.torrentFile.PieceLength
		boundaries := p.torrentFile.FileBoundariesMapping

		for {
			select {
			case <- ctx.Done():
				return
			case update := <- inputChan:
				logrus.Debugf("Got priority update in StartPriorityUpdating: %v", update)
				for _, file := range boundaries {
					if file.FileName == update.FileName {
						if update.ByteIdx > (file.End - file.Start) {
							logrus.Errorf("Wow! update.ByteIdx (%v) > (file.End - file.Start) (%v)", update.ByteIdx, file.End - file.Start)
							continue
						}
						overAllByteIdx := file.Start + update.ByteIdx
						pieceIdx := int(overAllByteIdx) / pLen
						if pieceIdx > nPieces{
							logrus.Errorf("Wow! pieceIdx (%v) > int64(len(p.torrentFile.PieceHashes)) (%v)", pieceIdx, nPieces)
							continue
						}
						outputChan <- pieceIdx
					}
				}
			}
		}
	}()

	return outputChan
}
