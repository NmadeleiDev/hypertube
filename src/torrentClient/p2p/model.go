package p2p

import (
	"torrentClient/client"
	"torrentClient/loadMaster"
)

const MaxBlockSize = 16384

const MaxBacklog = 5

type TorrentMeta struct {
	ActivatedClientsChan     <- chan *client.Client
	DeadPeersChan            chan <- *client.Client
	PieceLoadPriorityUpdates <- chan int
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
