package torrentfile

import (
	"crypto/md5"
	"fmt"
	"strings"
	"sync"
	"time"
)

type TorrentFile struct {
	mu sync.Mutex

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

func (t *TorrentFile) GetFileBoundariesMapping() (res []FileBoundaries) {
	t.mu.Lock()
	res = make([]FileBoundaries, len(t.FileBoundariesMapping))
	copy(res, t.FileBoundariesMapping)
	t.mu.Unlock()
	return res
}

func (t *TorrentFile) SetFileBoundariesMapping(src []FileBoundaries) {
	t.mu.Lock()
	t.FileBoundariesMapping = make([]FileBoundaries, len(src))
	copy(t.FileBoundariesMapping, src)
	t.mu.Unlock()
}

func (t *TorrentFile) GetFiles() (res []bencodeTorrentFile) {
	t.mu.Lock()
	res = make([]bencodeTorrentFile, len(t.Files))
	copy(res, t.Files)
	t.mu.Unlock()
	return res
}

func (t *TorrentFile) SetFiles(src []bencodeTorrentFile) {
	t.mu.Lock()
	t.Files = make([]bencodeTorrentFile, len(src))
	copy(t.Files, src)
	t.mu.Unlock()
}

func (t *TorrentFile) GetName() (res string) {
	t.mu.Lock()
	res = t.Name
	t.mu.Unlock()
	return res
}

func (t *TorrentFile) GetPieceLength() (res int) {
	t.mu.Lock()
	res = t.PieceLength
	t.mu.Unlock()
	return res
}

func (t *TorrentFile) GetLength() (res int) {
	t.mu.Lock()
	res = t.Length
	t.mu.Unlock()
	return res
}

func (t *TorrentFile) GetFileId() (res string) {
	t.mu.Lock()
	res = t.SysInfo.FileId
	t.mu.Unlock()
	return res
}

func (t *TorrentFile) GetInfoHash() (res [20]byte) {
	t.mu.Lock()
	copy(res[:], t.InfoHash[:])
	t.mu.Unlock()
	return res
}

func (t *TorrentFile) SetInfoHash(src [20]byte) {
	t.mu.Lock()
	copy(t.InfoHash[:], src[:])
	t.mu.Unlock()
}

func (t *TorrentFile) GetPieceHashes() (res [][20]byte) {
	t.mu.Lock()
	res = make([][20]byte, len(t.PieceHashes))
	copy(res, t.PieceHashes)
	t.mu.Unlock()
	return res
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

	mu      sync.Mutex
	isValid bool
}

func (u *UdpConnManager) SetValid(val bool) {
	u.mu.Lock()
	u.isValid = val
	u.mu.Unlock()
}

func (u *UdpConnManager) IsValid() bool {
	u.mu.Lock()
	val := u.isValid
	u.mu.Unlock()

	return val
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

func (b *bencodeTorrentFile) Extension() string {
	if len(b.Path) == 0 {
		return ""
	}
	lastPart := b.Path[len(b.Path) - 1]
	split := strings.Split(lastPart, ".")
	if len(split) < 2 {
		return ""
	}
	return split[len(split) - 1]
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

