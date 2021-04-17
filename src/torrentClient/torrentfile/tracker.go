package torrentfile

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"torrentClient/peers"

	"github.com/jackpal/bencode-go"
	"github.com/sirupsen/logrus"
)

type bencodeTrackerResp struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

const (
	protocolId = 0x41727101980
	connectAction = 0
	announceAction = 1
	scrapeAction = 2
)

type Tracker struct {
	Announce		string
	TransactionId	uint32
	ConnectionId	uint64
	MyPeerId		[20]byte
	MyPeerPort		uint16
	TrackerCallInterval		time.Duration
	UdpManager	*UdpConnManager

	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
}

//func (t *TorrentFile) RequestPeers() ([]peers.Peer, error) {
//	allPeers := make([]peers.Peer, 0, 50)
//	if peerIds, err := t.CallFittingScheme(t.Announce); err == nil {
//		allPeers = append(allPeers, peerIds...)
//	} else {
//		logrus.Errorf("Error calling main announce: %v", err)
//	}
//
//	for _, announce := range t.AnnounceList {
//		if peerIds, err := t.CallFittingScheme(announce); err == nil {
//			allPeers = append(allPeers, peerIds...)
//		} else {
//			logrus.Errorf("Error calling announce list member: %v", err)
//		}
//	}
//	if len(allPeers) > 0 {
//		return allPeers, nil
//	} else {
//		return nil, fmt.Errorf("failed to call any tracker")
//	}
//}

// TODO переписать под многопоточность: не использовать глобальную структуру (готово)
func (t *Tracker) CallFittingScheme() ([]peers.Peer, error) {
	trackerUrl, err := url.Parse(t.Announce)
	if err != nil {
		logrus.Errorf("Error parse tracker url: %v", err)
		return nil, err
	}

	if trackerUrl.Scheme == "http" {
		return t.callHttpTracker()
	} else if trackerUrl.Scheme == "udp" {
		return t.callUdpTracker()
	} else {
		return nil, fmt.Errorf("unsupported url scheme: %v; url: %v", trackerUrl.Scheme, t.Announce)
	}
}

func (t *Tracker) callUdpTracker() ([]peers.Peer, error) {
	trackerUrl, err := url.Parse(t.Announce)
	if err != nil {
		logrus.Errorf("Error parsing tracker url (%v): %v", t.Announce, err)
		return nil, err
	}
	t.UdpManager, err = OpenUdpSocket(trackerUrl)
	if err != nil {
		return nil, err
	}

	defer func() {
		t.UdpManager.ExitChan <- 1
	}()

	if err := t.makeConnectUdpReq(); err != nil {
		return nil, err
	}

	t.makeScrapeUdpReq()

	parsedPeers, err := t.makeAnnounceUdpReq()
	return parsedPeers, err
}

func (t *Tracker) callHttpTracker() ([]peers.Peer, error) {
	urlStr, err := t.buildHttpTrackerURL(t.Announce)
	if err != nil {
		return nil, err
	}
	c := &http.Client{Timeout: 15 * time.Second}
	resp, err := c.Get(urlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to send GET with client: %v; url: %v", err, urlStr)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading resp body: %v", err)
	}

	trackerResp := bencodeTrackerResp{}
	err = bencode.Unmarshal(bytes.NewBuffer(body), &trackerResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshal bencode: %v; body: %v", err, string(body))
	}
	return peers.Unmarshal([]byte(trackerResp.Peers))
}

func (t *Tracker) makeConnectUdpReq() error {
	req, err := t.buildUdpTrackerConnectReq()
	if err != nil {
		return err
	}

	t.UdpManager.Send <- req

	var body []byte

	timer := time.NewTimer(time.Second * 3)
	select {
	case <- timer.C:
		return fmt.Errorf("tracker call timed out")
	case data := <- t.UdpManager.Receive:
		body = data
		timer.Stop()
	}

	transId := binary.BigEndian.Uint32(body[4:9])
	if transId != t.TransactionId {
		logrus.Errorf("Tracker resp trans_id (%v) != saved trans_id (%v)", transId, t.TransactionId)
		// выйти?
	}
	t.ConnectionId = binary.BigEndian.Uint64(body[8:])

	logrus.Infof("Connect announce resp: conn_id=%v action=%v trans_id=%v", t.ConnectionId, binary.BigEndian.Uint32(body[:4]), binary.BigEndian.Uint32(body[4:8]))
	return nil
}

func (t *Tracker) makeAnnounceUdpReq() ([]peers.Peer, error) {
	req, err := t.buildUdpTrackerAnnounceReq()
	if err != nil {
		return nil, err
	}

	t.UdpManager.Send <- req
	var body []byte

	timer := time.NewTimer(time.Second * 10)

	select {
	case <- timer.C:
		return nil, fmt.Errorf("tracker call timed out")
	case data := <- t.UdpManager.Receive:
		body = data
		timer.Stop()
	}

	transId := binary.BigEndian.Uint32(body[4:8])
	if transId != t.TransactionId {
		logrus.Errorf("Tracker resp trans id (%v) != saved trans id (%v)", transId, t.TransactionId)
		// выйти?
	}
	interval := binary.BigEndian.Uint32(body[8:12])
	leechers := binary.BigEndian.Uint32(body[12:16])
	seeders := binary.BigEndian.Uint32(body[16:20])

	logrus.Infof("Interval = %v; leechers = %v; seeders = %v;", interval, leechers, seeders)
	parsedPeers, err := peers.Unmarshal(body[20:])
	logrus.Infof("Got peers: %v", parsedPeers)
	t.TrackerCallInterval = time.Duration(interval)
	return parsedPeers, err
}

func (t *Tracker) makeScrapeUdpReq() {
	req, err := t.buildScrapeUdpReq()
	if err != nil {
		logrus.Errorf("Error building scrape req: %v", err)
		return
	}

	t.UdpManager.Send <- req
	var body []byte

	timer := time.NewTimer(time.Second * 10)

	select {
	case <- timer.C:
		logrus.Errorf("tracker call timed out")
		return
	case data := <- t.UdpManager.Receive:
		body = data
		timer.Stop()
	}

	transId := binary.BigEndian.Uint32(body[4:8])
	if transId != t.TransactionId {
		logrus.Errorf("Tracker resp trans id (%v) != saved trans id (%v)", transId, t.TransactionId)
		// выйти?
	}
	seeders := binary.BigEndian.Uint32(body[8:12])
	completed := binary.BigEndian.Uint32(body[12:16])
	leechers := binary.BigEndian.Uint32(body[16:20])

	logrus.Infof("Scrape res: completed = %v; leechers = %v; seeders = %v;", completed, leechers, seeders)
	return
}