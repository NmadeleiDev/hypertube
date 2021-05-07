package torrentfile

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/jackpal/bencode-go"
	"github.com/sirupsen/logrus"
)

func (t *torrentsManager) ParseReaderToTorrent(body io.Reader) (*TorrentFile, error) {
	btoSingle := bencodeTorrentSingleFile{}
	btoMultiFile := bencodeTorrentMultiFiles{}

	readBody, err := io.ReadAll(body)
	if err != nil {
		logrus.Errorf("error readall body: %v", err)
		return nil, err
	}
	err = bencode.Unmarshal(bytes.NewBuffer(readBody), &btoSingle)
	if err != nil {
		return nil, err
	}
	err = bencode.Unmarshal(bytes.NewBuffer(readBody), &btoMultiFile)
	if err != nil {
		return nil, err
	}
	logrus.Infof("Parsed torrent!")

	var result *TorrentFile

	if btoSingle.Info.Length == 0 {
		result, err = btoMultiFile.toTorrentFile()
	} else {
		result, err = btoSingle.toTorrentFile()
	}
	if err != nil {
		logrus.Errorf("Error creating torret from bto: %v", err)
		return nil, err
	}
	if result == nil {
		logrus.Errorf("Parsed torrent is nil for: %v", readBody)
		return nil, fmt.Errorf("failed to parse torrent")
	}

	logrus.Infof("Bto info: name='%v'; len=%v; files = %v; pieces = (%v)",
		result.Name, result.Length, result.Files,
		[]string{hex.EncodeToString(result.PieceHashes[0][:]),
			hex.EncodeToString(result.PieceHashes[1][:])})
	return result, nil
}


