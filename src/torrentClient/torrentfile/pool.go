package torrentfile

import (
	"context"
	"time"

	"torrentClient/client"
	"torrentClient/peers"

	"github.com/sirupsen/logrus"
)

const (
	reconnectAttemptWait = time.Second * 30
	reconnectAttempts = 5
)

func (p *PeersPool) StartRefreshing(ctx context.Context)  {
	if p.ClientFactoryChan == nil {
		logrus.Errorf("Pool not initialized!")
		return
	}

	announceList := make([]string, len(p.torrent.AnnounceList) + 1)
	announceList[0] = p.torrent.Announce
	copy(announceList[1:], p.torrent.AnnounceList)

	sentPeersMap := make(map[string]bool, 50)

	for _, announce := range announceList {
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
		go func(ctx context.Context) {
			timer := time.NewTimer(time.Second)
			for {
				select {
				case <- ctx.Done():
					return
				case <- timer.C:
					rawPeers, err := tracker.CallFittingScheme()
					if err != nil {
						logrus.Errorf("Error requesting peers: %v", err)
						return
					}

					for _, peer := range rawPeers {
						if isSet, exists := sentPeersMap[peer.GetAddr()]; exists && isSet {
							continue
						}

						go func(peerToInit peers.Peer) {
							activeClient := p.InitPeer(&peerToInit)
							if activeClient != nil {
								sentPeersMap[peerToInit.GetAddr()] = true
								peerToInit.IsDead = false
								p.ClientFactoryChan <- activeClient
								logrus.Infof("Wrote peer %v to active clients chan", activeClient.GetShortInfo())
							} else {
								peerToInit.IsDead = true
								go p.StartConnAttempts(ctx, peerToInit)
							}
						}(peer)
					}
					timer.Reset(time.Second * tracker.TrackerCallInterval)
				}
			}
		}(ctx)
	}

	p.ListenForDeadPeers(ctx)
}

func (p *PeersPool) ListenForDeadPeers(ctx context.Context) {
	for {
		select {
		case <- ctx.Done():
			return
		case peer := <- p.ClientFactoryChan: // сюда из загрузчика приходят пиры, с которыми оборвалось соединение
			logrus.Debugf("Got dead peer %v, tring to raise him...", peer.GetShortInfo())
			go p.StartConnAttempts(ctx, peer.GetPeer())
		}
	}
}

func (p *PeersPool) StartConnAttempts(ctx context.Context, peer peers.Peer) {
	nAttempt := 1

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <- ctx.Done():
			return
		case <- ticker.C:
			activeClient := p.InitPeer(&peer)
			if activeClient != nil {
				peer.IsDead = false
				p.ClientFactoryChan <- activeClient
				logrus.Infof("Raised peer %v from %v attempts, sending to chan", activeClient.GetShortInfo(), nAttempt)
				return
			} else {
				peer.IsDead = true
			}
			ticker.Reset(time.Second * time.Duration(nAttempt * 5))
			nAttempt ++
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
	p.ClientFactoryChan = make(chan *client.Client, 10)
}

func (p *PeersPool) DestroyPool()  {
	close(p.ClientFactoryChan)
}
