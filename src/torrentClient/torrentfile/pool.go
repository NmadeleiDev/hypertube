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

var generalTrackerList = []string{"udp://tracker.torrent.eu.org:451/announce", "udp://tracker.doko.moe:6969/announce", "udp://thetracker.org:80/announce", "udp://santost12.xyz:6969/announce", "udp://bt.xxx-tracker.com:2710/announce", "udp://tracker.vanitycore.co:6969/announce", "udp://zephir.monocul.us:6969/announce", "http://grifon.info:80/announce", "udp://bt.aoeex.com:8000/announce", "udp://tracker.tiny-vps.com:6969/announce", "udp://tracker.tvunderground.org.ru:3218/announce", "udp://tracker.halfchub.club:6969/announce", "udp://retracker.nts.su:2710/announce", "udp://wambo.club:1337/announce", "udp://tracker.dutchtracking.com:6969/announce", "udp://tc.animereactor.ru:8082/announce", "udp://tracker.justseed.it:1337/announce", "udp://tracker.leechers-paradise.org:6969/announce", "udp://tracker.opentrackr.org:1337/announce", "https://open.kickasstracker.com:443/announce", "udp://tracker.coppersurfer.tk:6969/announce", "udp://open.stealth.si:80/announce", "http://retracker.mgts.by:80/announce", "udp://inferno.demonoid.pw:3418/announce", "udp://tracker.cypherpunks.ru:6969/announce", "udp://tracker.grepler.com:6969/announce", "udp://tracker.flashtorrents.org:6969/announce", "udp://tracker.yoshi210.com:6969/announce", "udp://tracker.tiny-vps.com:6969/announce", "udp://tracker.internetwarriors.net:1337/announce", "udp://mgtracker.org:2710/announce", "http://tracker.yoshi210.com:6969/announce", "http://tracker.tiny-vps.com:6969/announce", "udp://tracker.filetracker.pl:8089/announce", "udp://tracker.ex.ua:80/announce", "http://mgtracker.org:2710/announce", "udp://tracker.aletorrenty.pl:2710/announce", "http://tracker.filetracker.pl:8089/announce", "http://tracker.ex.ua/announce", "http://mgtracker.org:6969/announce", "http://retracker.krs-ix.ru:80/announce", "udp://tracker2.indowebster.com:6969/announce", "http://thetracker.org:80/announce", "http://tracker.bittor.pw:1337/announce", "udp://tracker.kicks-ass.net:80/announce", "udp://tracker.aletorrenty.pl:2710/announce", "http://tracker.aletorrenty.pl:2710/announce", "http://tracker.bittorrent.am/announce", "udp://tracker.kicks-ass.net:80/announce", "http://tracker.kicks-ass.net/announce", "http://tracker.baravik.org:6970/announce", "http://tracker.dutchtracking.com/announce", "http://tracker.dutchtracking.com:80/announce", "udp://tracker4.piratux.com:6969/announce", "http://tracker.internetwarriors.net:1337/announce", "udp://tracker.skyts.net:6969/announce", "http://tracker.dutchtracking.nl/announce", "http://tracker2.itzmx.com:6961/announce", "http://tracker2.wasabii.com.tw:6969/announce", "http://www.wareztorrent.com:80/announce", "udp://bt.xxx-tracker.com:2710/announce", "udp://tracker.eddie4.nl:6969/announce", "udp://tracker.grepler.com:6969/announce", "udp://tracker.mg64.net:2710/announce", "udp://tracker.flashtorrents.org:6969/announce", "http://tracker.tfile.me/announce", "http://tracker1.wasabii.com.tw:6969/announce", "udp://tracker.bittor.pw:1337/announce", "http://tracker.tvunderground.org.ru:3218/announce", "http://tracker.grepler.com:6969/announce", "http://tracker.flashtorrents.org:6969/announce", "http://retracker.gorcomnet.ru/announce", "udp://tracker.sktorrent.net:6969/announce", "udp://tracker.sktorrent.net:6969", "udp://public.popcorn-tracker.org:6969/announce", "udp://tracker.ilibr.org:80/announce", "udp://tracker.kuroy.me:5944/announce", "udp://tracker.mg64.net:6969/announce", "udp://tracker.cyberia.is:6969/announce", "http://tracker.devil-torrents.pl:80/announce", "udp://tracker2.christianbro.pw:6969/announce", "udp://retracker.lanta-net.ru:2710/announce", "udp://tracker.internetwarriors.net:1337/announce", "udp://ulfbrueggemann.no-ip.org:6969/announce", "http://torrentsmd.eu:8080/announce", "udp://peerfect.org:6969/announce", "udp://tracker.swateam.org.uk:2710/announce", "http://ns349743.ip-91-121-106.eu:80/announce", "http://torrentsmd.me:8080/announce", "http://agusiq-torrents.pl:6969/announce", "http://fxtt.ru:80/announce", "udp://tracker.vanitycore.co:6969/announce", "udp://explodie.org:6969"}

type PeersPool struct {
	Peers                []*peers.Peer
	ClientMaker          PeersInitializer

	torrent *TorrentFile
}

type PeersInitializer struct {
	RawPeersChan	chan peers.Peer
	DeadPeersChan		chan *client.Client
	InitializedPeersChan	chan *client.Client
}

func (p *PeersPool) SetTorrent(src *TorrentFile) {
	p.torrent = src
}

func (p *PeersPool) StartRefreshing(ctx context.Context)  {
	go p.ClientMaker.ListenForRawPeers(ctx, p.torrent.Download.MyPeerId, p.torrent.InfoHash)
	go p.ClientMaker.ListenForDeadPeers(ctx)

	announceList := make([]string, len(p.torrent.AnnounceList) + 1 + len(generalTrackerList))
	announceList[0] = p.torrent.Announce
	copy(announceList[1:1 + len(p.torrent.AnnounceList)], p.torrent.AnnounceList)
	copy(announceList[1 + len(p.torrent.AnnounceList):], generalTrackerList)

	sentPeersMap := make(map[string]bool, 50)
	trackersCalled := make([]string, 0, len(announceList))

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
		go func(ctx context.Context, trackerInstance Tracker) {
			timer := time.NewTimer(time.Nanosecond)
			for {
				select {
				case <- ctx.Done():
					return
				case <- timer.C:
					rawPeers, err := trackerInstance.CallFittingScheme()
					trackersCalled = append(trackersCalled, trackerInstance.Announce)
					if err != nil {
						logrus.Errorf("Error requesting peers: %v", err)
						return
					}
					logrus.Debugf("Got peers from tracker(call int=%v) %v: %v", trackerInstance.TrackerCallInterval, trackerInstance.Announce, rawPeers)

					for _, peer := range rawPeers {
						if isSet, exists := sentPeersMap[peer.GetAddr()]; exists && isSet {
							continue
						}
						p.ClientMaker.RawPeersChan <- peer
						sentPeersMap[peer.GetAddr()] = true
					}
					if trackerInstance.TrackerCallInterval < time.Minute {
						trackerInstance.TrackerCallInterval = time.Minute
					}
					timer.Reset(time.Second * trackerInstance.TrackerCallInterval)
				}
			}
		}(ctx, tracker)
	}
	logrus.Debugf("")
}

func (pi *PeersInitializer) InitPeer(ctx context.Context, peer *peers.Peer, myId [20]byte, infoHash [20]byte) *client.Client {
	c, err := client.New(*peer, myId, infoHash)
	if err != nil {
		logrus.Errorf("Could not handshake with %s. Err: %v", peer.GetAddr(), err)
		peer.IsDead = true
		return nil
	}

	c.SendUnchoke()
	c.SendInterested()

	unchokeWaitCtx, waitCancel := context.WithTimeout(ctx, time.Second * 45)
	defer waitCancel()

	logrus.Debugf("Waiting for unchoke from %v", c.Peer.GetAddr())
	ok, err := c.WaitForUnchoke(unchokeWaitCtx)
	if err != nil {
		logrus.Errorf("Error waiting for unchoke: %v", err)
		c.Conn.Close()
		return nil
	}
	if !ok {
		c.Conn.Close()
		return nil
	}

	logrus.Infof("Completed handshake and got unchoke with %s", peer.GetAddr())
	c.Peer.IsDead = false
	return c
}

func (p *PeersPool) Init() {
	p.ClientMaker.Init()
}

func (pi *PeersInitializer) Init() {
	pi.RawPeersChan = make(chan peers.Peer, 400)
	pi.DeadPeersChan = make(chan *client.Client, 200)
	pi.InitializedPeersChan = make(chan *client.Client, 200)
}

func (pi *PeersInitializer) Destroy() {
	close(pi.RawPeersChan)
	close(pi.InitializedPeersChan)
}

func (pi *PeersInitializer) ListenForRawPeers(ctx context.Context, myId [20]byte, infoHash [20]byte) {
	peersInProgress := 0
	maxJobs := 200
	jobs := make(chan struct{}, maxJobs)
	defer close(jobs)

	for {
		select {
		case <- ctx.Done():
			return
		case rawPeer := <- pi.RawPeersChan:
			peerToInit := rawPeer
			peersInProgress ++

			jobs <- struct{}{}
			logrus.Debugf("Trying to init peer %v, peersInProgress=%v", rawPeer.GetAddr(), peersInProgress)
			go func(peer peers.Peer) {
				defer func() {
					if r := recover(); r != nil {
						logrus.Debugf("Recovered in peer init: %v", r)
					}
					<- jobs
					peersInProgress --
				}()

				activeClient := pi.InitPeer(ctx, &peer, myId, infoHash)
				if activeClient != nil {
					peer.IsDead = false
					pi.InitializedPeersChan <- activeClient
					logrus.Infof("Wrote peer %v to active clients chan; peers in progress=%v", activeClient.GetShortInfo(), peersInProgress)
				} else {
					logrus.Debugf("Failed to init peer %v, peersInProgress=%v", peer.GetAddr(), peersInProgress)
					peer.IsDead = true
					pi.DeadPeersChan <- &client.Client{Peer: peer}
				}
			}(peerToInit)
		}
	}
}

func (pi *PeersInitializer) ListenForDeadPeers(ctx context.Context) {
	for {
		select {
		case <- ctx.Done():
			logrus.Debugf("Exiting ListenForDeadPeers cause ctx done")
			return
		case peer := <- pi.DeadPeersChan: // сюда из загрузчика приходят пиры, с которыми оборвалось соединение
			logrus.Debugf("Got dead peer %v, tring to raise him...", peer.GetShortInfo())
			pi.RawPeersChan <- peer.Peer
		}
	}
}

func (p *PeersPool) DestroyPool()  {
	p.ClientMaker.Destroy()
}
