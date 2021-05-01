package p2p

import (
	"bytes"
	"context"
	"crypto/sha1"
	"fmt"
	"time"

	"torrentClient/client"
	"torrentClient/db"
	"torrentClient/message"

	"github.com/sirupsen/logrus"
)


func (state *pieceProgress) readMessage() error {
	msg, err := state.client.Read() // this call blocks
	if err != nil {
		return err
	}

	if msg == nil { // keep-alive
		return nil
	}

	switch msg.ID {
	case message.MsgUnchoke:
		state.client.Choked = false
		logrus.Infof("Got UNCHOKE from %v", state.client.GetShortInfo())
	case message.MsgChoke:
		state.client.Choked = true
	case message.MsgHave:
		index, err := message.ParseHave(msg)
		if err != nil {
			return err
		}
		state.client.Bitfield.SetPiece(index)
	case message.MsgPiece:
		n, err := message.ParsePiece(state.index, state.buf, msg)
		if err != nil {
			return err
		}
		state.downloaded += n
		state.backlog--
	}
	return nil
}

func attemptDownloadPiece(c *client.Client, pw *pieceWork) ([]byte, error) {
	logrus.Debugf("Attempting to download piece (len=%v, idx=%v)", pw.length, pw.index)

	state := pieceProgress{
		index:  pw.index,
		client: c,
		buf:    make([]byte, pw.length),
	}

	//Setting a deadline helps get unresponsive peers unstuck.
	//30 seconds is more than enough time to download a 262 KB piece
	defer c.Conn.SetDeadline(time.Time{}) // Disable the deadline

	for state.downloaded < pw.length {
		c.Conn.SetDeadline(time.Now().Add(30 * time.Second))
		// If unchoked, send requests until we have enough unfulfilled requests
		if !state.client.Choked {
			logrus.Debugf("Downloading from %v. State: idx=%v, downloaded=%v (%v%%)", c.GetShortInfo(), state.index, state.downloaded, (state.downloaded * 100) / pw.length)
			for state.backlog < MaxBacklog && state.requested < pw.length {
				blockSize := MaxBlockSize
				// Last block might be shorter than the typical block
				if pw.length - state.requested < blockSize {
					blockSize = pw.length - state.requested
				}

				err := c.SendRequest(pw.index, state.requested, blockSize)
				if err != nil {
					return nil, err
				}
				state.backlog++
				state.requested += blockSize
			}
		} else {
			logrus.Warnf("CHOKED by %v for idx=%v, waiting for unchoke", state.client.GetShortInfo(), pw.index)
		}

		err := state.readMessage()
		if err != nil {
			return nil, fmt.Errorf("read msg err: %v", err)
		}
	}

	return state.buf, nil
}

func checkIntegrity(pw *pieceWork, buf []byte) error {
	hash := sha1.Sum(buf)
	if !bytes.Equal(hash[:], pw.hash[:]) {
		return fmt.Errorf("index %d failed integrity check", pw.index)
	}
	return nil
}

func (t *TorrentMeta) startDownloadWorker(c *client.Client, workQueue chan *pieceWork, results chan *pieceResult, deadPeerChan chan *client.Client) {
	defer c.Conn.Close()

	t.LoadStats.IncrActivePeers()

	for pw := range workQueue {
		if pw.length < 0 {
			logrus.Errorf("Attempting to download incorrect pw: %v", *pw)
			continue
		}

		if !c.Bitfield.HasPiece(pw.index) {
			workQueue <- pw // Put piece back on the queue
			continue
		}

		// Download the piece
		buf, err := attemptDownloadPiece(c, pw)
		if err != nil {
			logrus.Errorf("Throwing dead peer %v cause err: %v", c.GetShortInfo(), err)
			deadPeerChan <- c
			workQueue <- pw // Put piece back on the queue
			t.LoadStats.DecrActivePeers()
			return
		}

		err = checkIntegrity(pw, buf)
		if err != nil {
			logrus.Errorf("Check err: %v", err)
			workQueue <- pw // Put piece back on the queue
			continue
		}

		c.SendHave(pw.index)
		results <- &pieceResult{pw.index, buf}
	}
}

func (t *TorrentMeta) calculateBoundsForPiece(index int) (begin int, end int) {
	begin = index * t.PieceLength
	end = begin + t.PieceLength
	if end > t.Length {
		end = t.Length
	}
	return begin, end
}

func (t *TorrentMeta) calculatePieceSize(index int) int {
	begin, end := t.calculateBoundsForPiece(index)
	return end - begin
}

func (t *TorrentMeta) Download(ctx context.Context) error {
	logrus.Infof("starting download %v parts, file.len=%v, p.length=%v for %v",
		len(t.PieceHashes), t.Length, t.PieceLength, t.Name)

	loadedIdxs := db.GetFilesManagerDb().GetLoadedIndexesForFile(t.FileId)
	logrus.Debugf("Got loaded idxs: %v", loadedIdxs)

	priorityManager := prioritySorter{
		Pieces: make([]*pieceWork, 0, len(t.PieceHashes)),
		PriorityUpdates: t.PieceLoadPriorityUpdates,
	}

	for index, hash := range t.PieceHashes {
		if IntArrayContain(loadedIdxs, index) {
			if buf, start, size, ok := db.GetFilesManagerDb().GetPartDataByIdx(t.FileId, index); ok {
				t.ResultsChan <- LoadedPiece{Data: buf, Len: size, StartByte: start}
				t.LoadStats.IncrDone()
				continue
			} else {
				db.GetFilesManagerDb().DropDataPartByIdx(t.FileId, index)
			}
		}
		length := t.calculatePieceSize(index)
		//logrus.Debugf("Prepared piece idx=%v, len=%v", index, length)
		piece := pieceWork{index, hash, length}
		priorityManager.Pieces = append(priorityManager.Pieces, &piece)
	}

	if len(priorityManager.Pieces) == 0 { // значит все уже загружено
		return nil
	}

	results := make(chan *pieceResult, len(priorityManager.Pieces))

	defer close(results)

	topPriorityPieceChan := priorityManager.InitSorter(ctx)

	// Start workers as they arrive from Pool
	go func() {
		for {
			select {
			case <- ctx.Done():
				return
			case activeClient := <- t.ClientFactoryChan:
				if activeClient == nil {
					logrus.Errorf("Got nil active client...")
					continue
				}
				logrus.Infof("Got activated client: %v", activeClient.GetShortInfo())
				go t.startDownloadWorker(activeClient, topPriorityPieceChan, results, t.ClientFactoryChan)
			}
		}
	}()

	for t.LoadStats.CountDone() < len(t.PieceHashes) {
		select {
		case <- ctx.Done():
			logrus.Debugf("Got DONE in Download, exiting")
			return fmt.Errorf("load terminated by context")
		case res := <- results:
			if res == nil {
				logrus.Errorf("Piece result invalid: %v", res)
				continue
			}

			begin, end := t.calculateBoundsForPiece(res.index)
			t.ResultsChan <- LoadedPiece{Data: res.buf, Len: int64(end-begin), StartByte: int64(begin)}
			db.GetFilesManagerDb().SaveFilePart(t.FileId, res.buf, int64(begin), int64(end-begin), int64(res.index))
			//db.GetLoadedStateDb().AnnounceLoadedPart(t.FileId, fmt.Sprint(res.index), int64(begin), int64(end-begin))
			//db.GetLoadedStateDb().SaveLoadedPartInfo(t.FileId, fmt.Sprint(res.index), int64(begin), int64(end-begin))
			t.LoadStats.IncrDone()

			percent := t.LoadStats.GetLoadedPercent()
			logrus.Infof("(%d%%) Downloaded piece idx=%d from %v peers\n", percent, res.index, "n=unknown")
		}
	}
	return nil
}