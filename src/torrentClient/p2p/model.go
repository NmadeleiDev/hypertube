package p2p

import (
	"torrentClient/client"
	"torrentClient/peers"
)

// MaxBlockSize is the largest number of bytes a request can ask for
const MaxBlockSize = 16384

// MaxBacklog is the number of unfulfilled requests a client can have in its pipeline
const MaxBacklog = 5

// TorrentMeta holds data required to download a torrent from a list of peers
type TorrentMeta struct {
	ActiveClientsChan	<- chan *client.Client
	Peers       []peers.Peer
	PeerID      [20]byte
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
	FileId		string
}

func (t *TorrentMeta) CountWorkingPeers() (res int) {
	for _, peer := range t.Peers {
		if !peer.IsDead {
			res += 1
		}
	}
	return res
}

type pieceWork struct {
	index  int
	hash   [20]byte
	length int
}

type pieceResult struct {
	index int
	buf   []byte
}

type pieceProgress struct {
	index      int
	client     *client.Client
	buf        []byte
	downloaded int
	requested  int
	backlog    int
}
