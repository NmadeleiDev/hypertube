package torrentfile

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"torrent_client/db"
	"torrent_client/p2p"
	"torrent_client/parser/env"

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

func (t *torrentsManager) ReadTorrentFileFromHttpBody(body io.Reader) (TorrentFile, error) {
	return t.ParseReaderToTorrent(body)
}

func (t *torrentsManager) ParseReaderToTorrent(body io.Reader) (TorrentFile, error) {
	bto := bencodeTorrent{}
	err := bencode.Unmarshal(body, &bto)
	if err != nil {
		return TorrentFile{}, err
	} else {
		logrus.Infof("Parsed torrent!")
	}
	return bto.toTorrentFile()
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

	peers, err := t.requestPeers(peerID, Port)
	if err != nil {
		return fmt.Errorf("peers request error: %v", err)
	}

	torrent := p2p.Torrent{
		Peers:       peers,
		PeerID:      peerID,
		InfoHash:    t.InfoHash,
		PieceHashes: t.PieceHashes,
		PieceLength: t.PieceLength,
		Length:      t.Length,
		Name:        t.Name,
		FileId: 	 t.FileId,
	}

	db.GetFilesManagerDb().PreparePlaceForFile(torrent.FileId)
	logrus.Infof("Prepared table for parts,  starting download")

	err = torrent.Download()
	if err != nil {
		return fmt.Errorf("file download error: %e", err)
	}

	outFile, err := ioutil.TempFile(env.GetParser().GetFilesDir(), "loaded_*")
	if err != nil {
		return fmt.Errorf("create tempfile error: %e", err)
	}
	defer outFile.Close()

	db.GetFilesManagerDb().SaveFilePartsToFile(outFile, t.FileId)

	db.GetFilesManagerDb().SaveFileNameForReadyFile(t.FileId, outFile.Name())

	db.GetFilesManagerDb().RemoveFilePartsPlace(torrent.FileId)

	return nil
}

func (i *bencodeInfo) hash() ([20]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *i)
	if err != nil {
		return [20]byte{}, err
	}
	h := sha1.Sum(buf.Bytes())
	return h, nil
}

func (i *bencodeInfo) splitPieceHashes() ([][20]byte, error) {
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

func (bto *bencodeTorrent) toTorrentFile() (TorrentFile, error) {
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
		InfoHash:    infoHash,
		PieceHashes: pieceHashes,
		PieceLength: bto.Info.PieceLength,
		Length:      bto.Info.Length,
		Name:        bto.Info.Name,
	}
	return t, nil
}
