package torrentfile

import (
	"crypto/md5"
	"fmt"
	"strings"
	"time"

	"torrentClient/client"
	"torrentClient/peers"
)

// Port to listen on
const Port uint16 = 6881

// TorrentFile encodes the metadata from a .torrent file
type TorrentFile struct {
	Announce    string
	AnnounceList	[]string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Files		[]bencodeTorrentFile
	FileBoundariesMapping	 []FileBoundaries
	Name        string
	SysInfo     SystemInfo
	Download    DownloadUtils
}

type FileBoundaries struct {
	FileName string
	Index	int
	Start	int64
	End		int64
}

type SystemInfo struct {
	FileId		string
}

type DownloadUtils struct {
	TransactionId	uint32
	ConnectionId	uint64
	MyPeerId		[20]byte
	MyPeerPort		uint16
	TrackerCallInterval		time.Duration
	UdpManager	*UdpConnManager
}

type UdpConnManager struct {
	Receive chan []byte
	Send chan []byte
	ExitChan chan byte
	IsValid	bool
}

type bencodeInfoSingleFile struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

type bencodeInfoMultiFiles struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Files      []bencodeTorrentFile    `bencode:"files"`
	Name        string `bencode:"name"`
}

type bencodeTorrentFile struct {
	Length int      `bencode:"length"`
	Path   []string `bencode:"path"`
}

func (b *bencodeTorrentFile) EncodeFileName() string {
	hash := md5.Sum([]byte(strings.Join(b.Path, "")))
	return fmt.Sprintf("%x", hash[:])
}

type bencodeTorrentSingleFile struct {
	Announce     string                `bencode:"announce"`
	AnnounceList [][]string            `bencode:"announce-list"`
	Info         bencodeInfoSingleFile `bencode:"info"`
}

type bencodeTorrentMultiFiles struct {
	Announce     string                `bencode:"announce"`
	AnnounceList [][]string            `bencode:"announce-list"`
	Info         bencodeInfoMultiFiles `bencode:"info"`
}

func (bto *bencodeTorrentMultiFiles) SumFilesLength() int {
	res := 0
	for _, file := range bto.Info.Files {
		res += file.Length
	}
	return res
}


type PeersPool struct {
	Peers             []*peers.Peer
	ActiveClientsChan chan *client.Client

	torrent *TorrentFile
}

func (p *PeersPool) SetTorrent(src *TorrentFile) {
	p.torrent = src
}

