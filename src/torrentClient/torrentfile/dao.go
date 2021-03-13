package torrentfile

import "io"

type TorrentFilesManager interface {
	ReadTorrentFileFromFS(path string) (TorrentFile, error)
	ReadTorrentFileFromBytes(body io.Reader) (TorrentFile, error)
}
