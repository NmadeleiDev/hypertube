package p2p

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"log"
	"runtime"

	"torrentClient/client"
	"torrentClient/db"
	"torrentClient/message"
	"torrentClient/peersPool"

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
	logrus.Infof("Attempting to download piece (len=%v, idx=%v)", pw.length, pw.index)

	if pw.length < 0 {
		logrus.Errorf("Attempting to download incorrect pw: %v", *pw)
		return nil, fmt.Errorf("incorrect pw")
	}

	state := pieceProgress{
		index:  pw.index,
		client: c,
		buf:    make([]byte, pw.length),
	}

	// Setting a deadline helps get unresponsive peers unstuck.
	// 30 seconds is more than enough time to download a 262 KB piece
	//c.Conn.SetDeadline(time.Now().Add(300 * time.Second))
	//defer c.Conn.SetDeadline(time.Time{}) // Disable the deadline

	for state.downloaded < pw.length {
		// If unchoked, send requests until we have enough unfulfilled requests
		if !state.client.Choked {
			logrus.Infof("Downloading from %v. State: idx=%v, downloaded=%v (%v%%)", c.GetShortInfo(), state.index, state.downloaded, (state.downloaded * 100) / pw.length)
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
	//c, err := client.New(peer, t.PeerID, t.InfoHash)
	//if err != nil {
	//	logrus.Errorf("Could not handshake with %s. Err: %v, my peer id: %v", peer.GetAddr(), err, t.PeerID)
	//	deadPeerChan <- &peer
	//	return
	//}
	//defer c.Conn.Close()
	////logrus.Infof("Completed handshake with %s", peer.GetAddr())
	////logrus.Infof("Completed handshake! Client info: %v", c.GetClientInfo())
	//
	//c.SendUnchoke()
	//c.SendInterested()
	//
	//logrus.Infof("Unchoke and interested sent: %v", c.GetClientInfo())

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

//[
//	{18.219.47.231 6881 false}
//	{88.161.9.10 6881 false}	to
//	{117.120.10.87 26132 false}
//	{222.214.196.53 15000 false}	to
//	{154.21.23.134 16072 false}
//	{117.181.199.235 56961 false}	to
//{117.80.161.50 16881 false}		+u
//	{114.253.97.248 51413 false}	to
//	{110.184.240.204 22223 false}	to
//{106.122.206.163 51413 false}	+u
//]


// Download downloads the torrent. This stores the entire file in memory.
func (t *TorrentMeta) Download(peersPool peersPool.PeersPool) error {
	logrus.Infof("starting download %v parts, file.len=%v, p.length=%v for %v",
		len(t.PieceHashes), t.Length, t.PieceLength, t.Name)

	peersPool.StartRefreshing()
	// Init queues for workers to retrieve work and send results
	workQueue := make(chan *pieceWork, len(t.PieceHashes))
	results := make(chan *pieceResult)
	for index, hash := range t.PieceHashes {
		length := t.calculatePieceSize(index)
		//logrus.Infof("Putting piece to queue: idx=%v, len=%v", index, length)
		workQueue <- &pieceWork{index, hash, length}
	}

	// Start workers
	go func() {
		for {
			activeClient := <- peersPool.NewPeersChan
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
		log.Printf("(%0.2f%%) Downloaded piece #%d from %d peers\n", percent, res.index, numWorkers)
	}
	return nil
}