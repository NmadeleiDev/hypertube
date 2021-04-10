package peersPool

import (
	"torrentClient/client"
	"torrentClient/peers"
	"torrentClient/torrentfile"
)

type PeersPool struct {
	Peers	[]*peers.Peer
	NewPeersChan	chan *client.Client

	torrent *torrentfile.TorrentFile
}

func (p *PeersPool) SetTorrent(src *torrentfile.TorrentFile) {
	p.torrent = src
}
