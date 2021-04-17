package torrentfile

import (
	"time"

	"torrentClient/client"
	"torrentClient/peers"

	"github.com/sirupsen/logrus"
)

func (p *PeersPool) StartRefreshing()  {
	if p.ActiveClientsChan == nil {
		logrus.Errorf("Pool not initialized!")
		return
	}

	announceList := make([]string, len(p.torrent.AnnounceList) + 1)
	announceList[0] = p.torrent.Announce
	copy(announceList[1:], p.torrent.AnnounceList)

	sentPeersMap := make(map[string]bool, 50)

	for _, announce := range announceList {
		//trackerAddr := announce
		tracker := Tracker{
			Announce: announce,
			TransactionId: 0,
			ConnectionId: 0,
			MyPeerId: p.torrent.Download.MyPeerId,
			MyPeerPort: p.torrent.Download.MyPeerPort,
			TrackerCallInterval: 0,
			UdpManager: nil,
			InfoHash: p.torrent.InfoHash,
			PieceHashes: p.torrent.PieceHashes,
			PieceLength: p.torrent.PieceLength,
			Length: p.torrent.Length,
		}
		go func() {
			timer := time.NewTimer(time.Second)
			for {
				<- timer.C
				rawPeers, err := tracker.CallFittingScheme()
				if err != nil {
					logrus.Errorf("Error requesting peers: %v", err)
					return
				}

				for _, peer := range rawPeers {
					if isSet, exists := sentPeersMap[peer.GetAddr()]; exists && isSet {
						continue
					}

					activeClient := p.InitPeer(&peer)
					if activeClient != nil {
						sentPeersMap[peer.GetAddr()] = true
						p.ActiveClientsChan <- activeClient
						logrus.Infof("Wrote peer %v to active clients chan", activeClient.GetShortInfo())
					} else {
						peer.IsDead = true
					}
				}
				timer.Reset(time.Second * tracker.TrackerCallInterval)
			}
		}()
	}
}

func (p *PeersPool) StartRetyingPeerConn(returnPeer chan <- *client.Client) {
	ticker := time.NewTicker(time.Second * 60)
	defer ticker.Stop()

	for {
		<- ticker.C

		for _, peer := range p.Peers {
			if peer.IsDead {
				activeClient := p.InitPeer(peer)
				if activeClient != nil {
					peer.IsDead = false
					returnPeer <- activeClient
				}
			}
		}
	}
}

func (p *PeersPool) InitPeer(peer *peers.Peer) *client.Client {
	c, err := client.New(*peer, p.torrent.Download.MyPeerId, p.torrent.InfoHash)
	if err != nil {
		logrus.Errorf("Could not handshake with %s. Err: %v", peer.GetAddr(), err)
		return nil
	}
	//defer c.Conn.Close()
	//logrus.Infof("Completed handshake! Client info: %v", c.GetClientInfo())

	c.SendUnchoke()
	c.SendInterested()

	logrus.Infof("Completed handshake with %s", peer.GetAddr())
	return c
}

func (p *PeersPool) InitPool() {
	p.ActiveClientsChan = make(chan *client.Client, 10)
}
