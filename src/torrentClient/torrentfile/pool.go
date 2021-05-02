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

var generalTrackerList = []string{"udp://tracker.torrent.eu.org:451/announce", "http://t.nyaatracker.com:80/announce", "udp://tracker.doko.moe:6969/announce", "http://asnet.pw:2710/announce", "udp://thetracker.org:80/announce", "http://tracker.tfile.co:80/announce", "http://pt.lax.mx:80/announce", "udp://santost12.xyz:6969/announce", "https://tracker.bt-hash.com:443/announce", "udp://bt.xxx-tracker.com:2710/announce", "udp://tracker.vanitycore.co:6969/announce", "udp://zephir.monocul.us:6969/announce", "http://grifon.info:80/announce", "http://retracker.spark-rostov.ru:80/announce", "http://tr.kxmp.cf:80/announce", "http://tracker.city9x.com:2710/announce", "udp://bt.aoeex.com:8000/announce", "http://tracker.tfile.me:80/announce", "udp://tracker.tiny-vps.com:6969/announce", "http://retracker.telecom.by:80/announce", "http://tracker.electro-torrent.pl:80/announce", "udp://tracker.tvunderground.org.ru:3218/announce", "udp://tracker.halfchub.club:6969/announce", "udp://retracker.nts.su:2710/announce", "udp://wambo.club:1337/announce", "udp://tracker.dutchtracking.com:6969/announce", "udp://tc.animereactor.ru:8082/announce", "udp://tracker.justseed.it:1337/announce", "udp://tracker.leechers-paradise.org:6969/announce", "udp://tracker.opentrackr.org:1337/announce", "https://open.kickasstracker.com:443/announce", "udp://tracker.coppersurfer.tk:6969/announce", "udp://open.stealth.si:80/announce", "http://retracker.mgts.by:80/announce", "http://retracker.bashtel.ru:80/announce", "udp://inferno.demonoid.pw:3418/announce", "udp://tracker.cypherpunks.ru:6969/announce", "http://tracker.calculate.ru:6969/announce", "udp://tracker.sktorrent.net:6969/announce", "udp://tracker.grepler.com:6969/announce", "udp://tracker.flashtorrents.org:6969/announce", "udp://tracker.yoshi210.com:6969/announce", "udp://tracker.tiny-vps.com:6969/announce", "udp://tracker.internetwarriors.net:1337/announce", "udp://mgtracker.org:2710/announce", "http://tracker.yoshi210.com:6969/announce", "http://tracker.tiny-vps.com:6969/announce", "udp://tracker.filetracker.pl:8089/announce", "udp://tracker.ex.ua:80/announce", "http://mgtracker.org:2710/announce", "udp://tracker.aletorrenty.pl:2710/announce", "http://tracker.filetracker.pl:8089/announce", "http://tracker.ex.ua/announce", "http://mgtracker.org:6969/announce", "http://retracker.krs-ix.ru:80/announce", "udp://tracker2.indowebster.com:6969/announce", "http://thetracker.org:80/announce", "http://tracker.bittor.pw:1337/announce", "udp://tracker.kicks-ass.net:80/announce", "udp://tracker.aletorrenty.pl:2710/announce", "http://tracker.aletorrenty.pl:2710/announce", "http://tracker.bittorrent.am/announce", "udp://tracker.kicks-ass.net:80/announce", "http://tracker.kicks-ass.net/announce", "http://tracker.baravik.org:6970/announce", "http://tracker.dutchtracking.com/announce", "http://tracker.dutchtracking.com:80/announce", "udp://tracker4.piratux.com:6969/announce", "http://tracker.internetwarriors.net:1337/announce", "udp://tracker.skyts.net:6969/announce", "http://tracker.dutchtracking.nl/announce", "http://tracker2.itzmx.com:6961/announce", "http://tracker2.wasabii.com.tw:6969/announce", "udp://tracker.sktorrent.net:6969/announce", "http://www.wareztorrent.com:80/announce", "udp://bt.xxx-tracker.com:2710/announce", "udp://tracker.eddie4.nl:6969/announce", "udp://tracker.grepler.com:6969/announce", "udp://tracker.mg64.net:2710/announce", "udp://tracker.coppersurfer.tk:6969/announce", "http://tracker.opentrackr.org:1337/announce", "http://tracker.dutchtracking.nl:80/announce", "http://tracker.edoardocolombo.eu:6969/announce", "http://tracker.ex.ua:80/announce", "http://tracker.kicks-ass.net:80/announce", "http://tracker.mg64.net:6881/announce", "udp://tracker.flashtorrents.org:6969/announce", "http://tracker.tfile.me/announce", "http://tracker1.wasabii.com.tw:6969/announce", "udp://tracker.bittor.pw:1337/announce", "http://tracker.tvunderground.org.ru:3218/announce", "http://tracker.grepler.com:6969/announce", "udp://tracker.bittor.pw:1337/announce", "http://tracker.flashtorrents.org:6969/announce", "http://retracker.gorcomnet.ru/announce", "udp://tracker.sktorrent.net:6969/announce", "udp://tracker.sktorrent.net:6969", "udp://public.popcorn-tracker.org:6969/announce", "udp://tracker.ilibr.org:80/announce", "udp://tracker.kuroy.me:5944/announce", "udp://tracker.mg64.net:6969/announce", "udp://tracker.cyberia.is:6969/announce", "http://tracker.devil-torrents.pl:80/announce", "udp://tracker2.christianbro.pw:6969/announce", "udp://retracker.lanta-net.ru:2710/announce", "udp://tracker.internetwarriors.net:1337/announce", "udp://ulfbrueggemann.no-ip.org:6969/announce", "http://torrentsmd.eu:8080/announce", "udp://peerfect.org:6969/announce", "udp://tracker.swateam.org.uk:2710/announce", "http://ns349743.ip-91-121-106.eu:80/announce", "http://torrentsmd.me:8080/announce", "http://agusiq-torrents.pl:6969/announce", "http://fxtt.ru:80/announce", "udp://tracker.vanitycore.co:6969/announce", "udp://explodie.org:6969"}

func (p *PeersPool) StartRefreshing(ctx context.Context)  {
	if p.ClientFactoryChan == nil {
		logrus.Errorf("Pool not initialized!")
		return
	}

	announceList := make([]string, len(p.torrent.AnnounceList) + 1 + len(generalTrackerList))
	announceList[0] = p.torrent.Announce
	copy(announceList[1:1 + len(p.torrent.AnnounceList)], p.torrent.AnnounceList)
	copy(announceList[1 + len(p.torrent.AnnounceList):], generalTrackerList)

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
								p.StartConnAttempts(ctx, peerToInit)
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

	timer := time.NewTimer(time.Second)
	defer timer.Stop()

	for {
		select {
		case <- ctx.Done():
			return
		case <- timer.C:
			if nAttempt > 5 {
				return
			}
			activeClient := p.InitPeer(&peer)
			if activeClient != nil {
				peer.IsDead = false
				p.ClientFactoryChan <- activeClient
				logrus.Infof("Raised peer %v from %v attempts, sending to chan", activeClient.GetShortInfo(), nAttempt)
				return
			} else {
				peer.IsDead = true
			}
			timer.Reset(time.Second * time.Duration(nAttempt * 5))
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
