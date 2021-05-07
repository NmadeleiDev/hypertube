package torrentfile

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"

	"torrentClient/db"
	"torrentClient/magnetToTorrent"

	"github.com/sirupsen/logrus"
)

func GetManager() TorrentFilesManager {
	return &torrentsManager{}
}

type TorrentFilesManager interface {
	ReadTorrentFileFromFS(path string) (*TorrentFile, error)
	ReadTorrentFileFromBytes(body io.Reader) (*TorrentFile, error)
	LoadTorrentFileFromDB(fileId string) (*TorrentFile, error)
}

func (t *torrentsManager) ReadTorrentFileFromFS(path string) (*TorrentFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return t.ParseReaderToTorrent(file)
}

func (t *torrentsManager) ReadTorrentFileFromBytes(body io.Reader) (*TorrentFile, error) {
	return t.ParseReaderToTorrent(body)
}

func (t *torrentsManager) LoadTorrentFileFromDB(fileId string) (*TorrentFile, error)  {
	torrentBytes, magnetLink, ok := db.GetFilesManagerDb().GetTorrentOrMagnetForByFileId(fileId)
	if !ok {
		return nil, fmt.Errorf("record not found")
	}

	doChangeAnnounce := false

	if (torrentBytes == nil || len(torrentBytes) == 0) && len(magnetLink) > 0 {
		torrentBytes = magnetToTorrent.ConvertMagnetToTorrent(magnetLink)
		doChangeAnnounce = true
	}

	torrent, err := t.ReadTorrentFileFromBytes(bytes.NewBuffer(torrentBytes))
	if err != nil || torrent == nil {
		return nil, fmt.Errorf("error reading torrent file: %v", err)
	}
	torrent.SysInfo.FileId = fileId

	if doChangeAnnounce {
		trackerUrl := t.GetTrackersFromMagnet(magnetLink)
		logrus.Infof("Tracker url: %v", trackerUrl)
		torrent.Announce = trackerUrl
	}

	return torrent, nil
}

func (t *torrentsManager) GetTrackersFromMagnet(magnet string) string {
	decoded, err := url.ParseQuery(magnet)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return decoded.Get("tr")
}
