package p2p

import (
	"torrentClient/client"
	"torrentClient/loadMaster"
)

// MaxBlockSize is the largest number of bytes a request can ask for
const MaxBlockSize = 16384

// MaxBacklog is the number of unfulfilled requests a client can have in its pipeline
const MaxBacklog = 5

// TorrentMeta holds data required to download a torrent from a list of peers
type TorrentMeta struct {
	ActivatedClientsChan     <- chan *client.Client
	DeadPeersChan            chan <- *client.Client
	PieceLoadPriorityUpdates <- chan int64
	PeerID                   [20]byte
	InfoHash                 [20]byte
	PieceHashes              [][20]byte
	PieceLength              int
	Length                   int
	Name                     string
	FileId                   string
	ResultsChan              chan LoadedPiece
	LoadStats				*loadMaster.LoadEntry
}

type LoadedPiece struct {
	StartByte	int64
	Len		int64
	Data	[]byte
}

type pieceWork struct {
	index  int
	hash   [20]byte
	length int
	progress *pieceProgress
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
