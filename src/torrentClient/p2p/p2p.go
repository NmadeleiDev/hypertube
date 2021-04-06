package p2p

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"log"
	"runtime"
	"time"

	"torrentClient/client"
	"torrentClient/db"
	"torrentClient/message"
	"torrentClient/peers"

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
	state := pieceProgress{
		index:  pw.index,
		client: c,
		buf:    make([]byte, pw.length),
	}

	// Setting a deadline helps get unresponsive peers unstuck.
	// 30 seconds is more than enough time to download a 262 KB piece
	c.Conn.SetDeadline(time.Now().Add(30 * time.Second))
	defer c.Conn.SetDeadline(time.Time{}) // Disable the deadline

	for state.downloaded < pw.length {
		// If unchoked, send requests until we have enough unfulfilled requests
		if !state.client.Choked {
			for state.backlog < MaxBacklog && state.requested < pw.length {
				blockSize := MaxBlockSize
				// Last block might be shorter than the typical block
				if pw.length-state.requested < blockSize {
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
		return fmt.Errorf("Index %d failed integrity check", pw.index)
	}
	return nil
}

func (t *Torrent) startDownloadWorker(peer peers.Peer, workQueue chan *pieceWork, results chan *pieceResult, deadPeerChan chan <- *peers.Peer) {
	c, err := client.New(peer, t.PeerID, t.InfoHash)
	if err != nil {
		logrus.Errorf("Could not handshake with %s. Err: %v, my peer id: %v", peer.String(), err, t.PeerID)
		deadPeerChan <- &peer
		return
	}
	defer c.Conn.Close()
	logrus.Infof("Completed handshake with %s", peer.String())

	c.SendUnchoke()
	c.SendInterested()

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
			logrus.Errorf("Piece #%d failed integrity check", pw.index)
			workQueue <- pw // Put piece back on the queue
			continue
		}

		c.SendHave(pw.index)
		results <- &pieceResult{pw.index, buf}
	}
}

func (t *Torrent) calculateBoundsForPiece(index int) (begin int, end int) {
	begin = index * t.PieceLength
	end = begin + t.PieceLength
	if end > t.Length {
		end = t.Length
	}
	return begin, end
}

func (t *Torrent) calculatePieceSize(index int) int {
	begin, end := t.calculateBoundsForPiece(index)
	return end - begin
}

// Download downloads the torrent. This stores the entire file in memory.
func (t *Torrent) Download() error {
	logrus.Infof("starting download %v parts from %v peers for %v", len(t.PieceHashes), len(t.Peers), t.Name)
	// Init queues for workers to retrieve work and send results
	workQueue := make(chan *pieceWork, len(t.PieceHashes))
	results := make(chan *pieceResult)
	for index, hash := range t.PieceHashes {
		length := t.calculatePieceSize(index)
		workQueue <- &pieceWork{index, hash, length}
	}

	deadPeersChan := make(chan *peers.Peer, 10)
	alarmAllDead := make(chan byte)
	go func() {
		count := 0
		for {
			peer := <- deadPeersChan
			count += 1
			logrus.Warnf("Peer %v dead. Total dead = %v", peer.String(), count)

			if count == len(t.Peers) {
				alarmAllDead <- 1
			}
		}
	}()

	// Start workers
	for _, peer := range t.Peers {
		go t.startDownloadWorker(peer, workQueue, results, deadPeersChan)
	}

	defer close(workQueue)
	// Collect results into a buffer until full
	//buf := make([]byte, t.Length)
	donePieces := 0
	for donePieces < len(t.PieceHashes) {
		var res *pieceResult
		select {
			case data := <-results:
				res = data
			case <- alarmAllDead:
				logrus.Infof("Failed to download file, as all peers are dead!")
				return fmt.Errorf("failed to download file, as all peers are dead")
		}
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
		log.Printf("(%0.2f%%) Downloaded piece #%d from %d peers\n", percent, res.index, numWorkers)
	}
	return nil
}