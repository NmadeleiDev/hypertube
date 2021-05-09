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

		nPieces := len(p.torrentFile.GetPieceHashes())
		pLen := p.torrentFile.GetPieceLength()
		boundaries := p.torrentFile.GetFileBoundariesMapping()

		for {
			select {
			case <- ctx.Done():
				return
			case update := <- inputChan:
				logrus.Debugf("Got priority update in StartPriorityUpdating: %v; files: %v", update, boundaries)
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
						logrus.Debugf("Sending priority update for file %v, piece %v", file.FileName, pieceIdx)
						outputChan <- pieceIdx
						logrus.Debugf("Priority update for file %v, piece %v sended successfully!", file.FileName, pieceIdx)
					}
				}
			}
		}
	}()

	return outputChan
}
