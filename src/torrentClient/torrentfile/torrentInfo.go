package torrentfile

import (
	"bytes"
	"crypto/sha1"
	"fmt"

	"github.com/jackpal/bencode-go"
)

func (bto *bencodeTorrentSingleFile) toTorrentFile() (*TorrentFile, error) {
	infoHash, err := bto.Info.hash()
	if err != nil {
		return nil, err
	}
	pieceHashes, err := bto.Info.splitPieceHashes()
	if err != nil {
		return nil, err
	}
	t := TorrentFile{
		Announce:    bto.Announce,
		AnnounceList: UnfoldArray(bto.AnnounceList),
		InfoHash:    infoHash,
		PieceHashes: pieceHashes,
		Files: 		[]bencodeTorrentFile{{Path: []string{bto.Info.Name}, Length: bto.Info.Length}},
		PieceLength: bto.Info.PieceLength,
		Length:      bto.Info.Length,
		Name:        bto.Info.Name,
		SysInfo:      SystemInfo{},
		Download:     DownloadUtils{},
	}
	return &t, nil
}

func (bto *bencodeTorrentMultiFiles) toTorrentFile() (*TorrentFile, error) {
	infoHash, err := bto.Info.hash()
	if err != nil {
		return nil, err
	}
	pieceHashes, err := bto.Info.splitPieceHashes()
	if err != nil {
		return nil, err
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
	return &t, nil
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
