package p2p

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"runtime"
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

	if pw.length < 0 {
		logrus.Errorf("Attempting to download incorrect pw: %v", *pw)
		return nil, fmt.Errorf("incorrect pw")
	}

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
		}

		err := state.readMessage()
		if err != nil {
			return nil, err
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

func (t *TorrentMeta) startDownloadWorker(c *client.Client, workQueue chan *pieceWork, results chan *pieceResult) {
	defer c.Conn.Close()

	for pw := range workQueue {
		if !c.Bitfield.HasPiece(pw.index) {
			workQueue <- pw // Put piece back on the queue
			continue
		}

		// Download the piece
		buf, err := attemptDownloadPiece(c, pw)
		if err != nil {
			logrus.Errorf("exiting: %v", err)
			workQueue <- pw // Put piece back on the queue
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

// Download downloads the torrent. This stores the entire file in memory.
func (t *TorrentMeta) Download() error {
	logrus.Infof("starting download %v parts, file.len=%v, p.length=%v for %v",
		len(t.PieceHashes), t.Length, t.PieceLength, t.Name)

	// Init queues for workers to retrieve work and send results
	workQueue := make(chan *pieceWork, len(t.PieceHashes))
	results := make(chan *pieceResult)

	loadedIdxs := db.GetFilesManagerDb().GetLoadedIndexesForFile(t.FileId)
	logrus.Debugf("Got loaded idxs: %v", loadedIdxs)

	for index, hash := range t.PieceHashes {
		if IntArrayContain(loadedIdxs, index) {
			continue
		}
		length := t.calculatePieceSize(index)

		logrus.Debugf("Prepared piece idx=%v, len=%v", index, length)
		workQueue <- &pieceWork{index, hash, length}
	}

	// Start workers
	go func() {
		for {
			activeClient := <- t.ActiveClientsChan
			logrus.Infof("Got activated client: %v", activeClient.GetShortInfo())
			go t.startDownloadWorker(activeClient, workQueue, results)
		}
	}()
	//for _, peer := range t.Peers {
	//	go t.startDownloadWorker(peer, workQueue, results)
	//}

	defer close(workQueue)
	// Collect results into a buffer until full
	//buf := make([]byte, t.Length)
	donePieces := 0
	for donePieces < len(t.PieceHashes) {
		var res *pieceResult
		res = <- results
		if res == nil {
			logrus.Errorf("Piece result invalid: %v", res)
			continue
		}

		begin, end := t.calculateBoundsForPiece(res.index)
		//copy(buf[begin:end], res.buf)

		db.GetFilesManagerDb().SaveFilePart(t.FileId, res.buf, int64(begin), int64(end-begin), int64(res.index))
		db.GetLoadedStateDb().AnnounceLoadedPart(t.FileId, fmt.Sprint(res.index), int64(begin), int64(end-begin))
		db.GetLoadedStateDb().SaveLoadedPartInfo(t.FileId, fmt.Sprint(res.index), int64(begin), int64(end-begin))

		donePieces++

		percent := float64(donePieces) / float64(len(t.PieceHashes)) * 100
		numWorkers := runtime.NumGoroutine() - 1 // subtract 1 for main thread
		logrus.Infof("(%0.2f%%) Downloaded piece idx=%d from %d peers\n", percent, res.index, numWorkers)
	}
	return nil
}