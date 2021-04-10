package torrentfile

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"torrentClient/db"
	"torrentClient/p2p"
	"torrentClient/parser/env"

	"github.com/jackpal/bencode-go"
	"github.com/sirupsen/logrus"
)

type torrentsManager struct {
}

func (t *torrentsManager) ReadTorrentFileFromFS(path string) (TorrentFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return TorrentFile{}, err
	}
	defer file.Close()

	return t.ParseReaderToTorrent(file)
}

func (t *torrentsManager) ReadTorrentFileFromBytes(body io.Reader) (TorrentFile, error) {
	return t.ParseReaderToTorrent(body)
}

func (t *torrentsManager) ParseReaderToTorrent(body io.Reader) (TorrentFile, error) {
	btoSingle := bencodeTorrentSingleFile{}
	btoMultiFile := bencodeTorrentMultiFiles{}

	readBody, err := io.ReadAll(body)
	if err != nil {
		logrus.Errorf("error readall body: %v", err)
		return TorrentFile{}, err
	}
	err = bencode.Unmarshal(bytes.NewBuffer(readBody), &btoSingle)
	if err != nil {
		return TorrentFile{}, err
	}
	err = bencode.Unmarshal(bytes.NewBuffer(readBody), &btoMultiFile)
	if err != nil {
		return TorrentFile{}, err
	}
	logrus.Infof("Parsed torrent!")

	var result TorrentFile

	if btoSingle.Info.Length == 0 {
		result, err = btoMultiFile.toTorrentFile()
	} else {
		result, err = btoSingle.toTorrentFile()
	}
	if err != nil {
		logrus.Errorf("Error creating torret from bto: %v", err)
	}

	logrus.Infof("Bto info: name='%v'; len=%v; files = %v; pieces = (%v)",
		result.Name, result.Length, result.Files,
		[]string{hex.EncodeToString(result.PieceHashes[0][:]),
			hex.EncodeToString(result.PieceHashes[1][:])})
	return result, nil
}

func GetManager() TorrentFilesManager {
	return &torrentsManager{}
}

// DownloadToFile downloads a torrent and writes it to a file
func (t *TorrentFile) DownloadToFile() error {
	var peerID [20]byte
	_, err := rand.Read(peerID[:])
	if err != nil{
		return fmt.Errorf("read rand error: %v", err)
	}
	t.Download.MyPeerId = peerID
	t.Download.MyPeerPort = env.GetParser().GetTorrentPeerPort()

	peersPoolObj := PeersPool{}
	peersPoolObj.InitPool()
	peersPoolObj.SetTorrent(t)
	go peersPoolObj.StartRefreshing()

	//peers, err := t.RequestPeers()
	//if err != nil {
	//	return fmt.Errorf("peers request error: %v", err)
	//}

	torrent := p2p.TorrentMeta{
		//Peers:       peers,
		ActiveClientsChan: peersPoolObj.ActiveClientsChan,
		PeerID:      t.Download.MyPeerId,
		InfoHash:    t.InfoHash,
		PieceHashes: t.PieceHashes,
		PieceLength: t.PieceLength,
		Length:      t.Length,
		Name:        t.Name,
		FileId: 	 t.SysInfo.FileId,
	}

	db.GetFilesManagerDb().PreparePlaceForFile(torrent.FileId)
	//defer db.GetFilesManagerDb().RemoveFilePartsPlace(torrent.FileId)
	logrus.Infof("Prepared table for parts, starting download")

	err = torrent.Download()
	if err != nil {
		return fmt.Errorf("file download error: %v", err)
	}

	outFile, err := ioutil.TempFile(env.GetParser().GetFilesDir(), "loaded_*")
	if err != nil {
		return fmt.Errorf("create tempfile error: %e", err)
	}
	defer func() {
		outFile.Close()
	}()

	db.GetFilesManagerDb().SaveFilePartsToFile(outFile, t.SysInfo.FileId)
	db.GetFilesManagerDb().SaveFileNameForReadyFile(t.SysInfo.FileId, outFile.Name())
	return nil
}

func (i *bencodeInfoSingleFile) hash() ([20]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *i)
	if err != nil {
		return [20]byte{}, err
	}
	h := sha1.Sum(buf.Bytes())
	return h, nil
}

func (i *bencodeInfoMultiFiles) hash() ([20]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *i)
	if err != nil {
		return [20]byte{}, err
	}
	h := sha1.Sum(buf.Bytes())
	return h, nil
}

func (i *bencodeInfoSingleFile) splitPieceHashes() ([][20]byte, error) {
	hashLen := 20 // Length of SHA-1 hash
	buf := []byte(i.Pieces)
	if len(buf)%hashLen != 0 {
		err := fmt.Errorf("Received malformed pieces of length %d", len(buf))
		return nil, err
	}
	numHashes := len(buf) / hashLen
	hashes := make([][20]byte, numHashes)

	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], buf[i*hashLen:(i+1)*hashLen])
	}
	return hashes, nil
}

func (i *bencodeInfoMultiFiles) splitPieceHashes() ([][20]byte, error) {
	hashLen := 20 // Length of SHA-1 hash
	buf := []byte(i.Pieces)
	if len(buf)%hashLen != 0 {
		err := fmt.Errorf("Received malformed pieces of length %d", len(buf))
		return nil, err
	}
	numHashes := len(buf) / hashLen
	hashes := make([][20]byte, numHashes)

	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], buf[i*hashLen:(i+1)*hashLen])
	}
	return hashes, nil
}

func (bto *bencodeTorrentSingleFile) toTorrentFile() (TorrentFile, error) {
	infoHash, err := bto.Info.hash()
	if err != nil {
		return TorrentFile{}, err
	}
	pieceHashes, err := bto.Info.splitPieceHashes()
	if err != nil {
		return TorrentFile{}, err
	}
	t := TorrentFile{
		Announce:    bto.Announce,
		AnnounceList: UnfoldArray(bto.AnnounceList),
		InfoHash:    infoHash,
		PieceHashes: pieceHashes,
		PieceLength: bto.Info.PieceLength,
		Length:      bto.Info.Length,
		Name:        bto.Info.Name,
	}
	return t, nil
}

func (bto *bencodeTorrentMultiFiles) toTorrentFile() (TorrentFile, error) {
	infoHash, err := bto.Info.hash()
	if err != nil {
		return TorrentFile{}, err
	}
	pieceHashes, err := bto.Info.splitPieceHashes()
	if err != nil {
		return TorrentFile{}, err
	}
	t := TorrentFile{
		Announce:     bto.Announce,
		AnnounceList: UnfoldArray(bto.AnnounceList),
		InfoHash:     infoHash,
		PieceHashes:  pieceHashes,
		PieceLength:  bto.Info.PieceLength,
		Length:       bto.SumFilesLength(),
		Files:        bto.Info.Files,
		Name:         bto.Info.Name,
		SysInfo:      SystemInfo{},
		Download:     DownloadUtils{},
	}
	return t, nil
}

func UnfoldArray(src [][]string) []string {
	res := make([]string, 0, len(src))

	for _, item := range src {
		res = append(res, item...)
	}

	return res
}
