package torrentfile

import (
	"bytes"
	"context"
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

func (t *Tracker) CallFittingScheme(ctx context.Context) ([]peers.Peer, error) {
	trackerUrl, err := url.Parse(t.Announce)
	if err != nil {
		logrus.Errorf("Error parse tracker url: %v", err)
		return nil, err
	}

	if trackerUrl.Scheme == "http" {
		return t.callHttpTracker()
	} else if trackerUrl.Scheme == "udp" {
		return t.callUdpTracker(ctx)
	} else {
		return nil, fmt.Errorf("unsupported url scheme: %v; url: %v", trackerUrl.Scheme, t.Announce)
	}
}

func (t *Tracker) callUdpTracker(ctx context.Context) ([]peers.Peer, error) {
	trackerUrl, err := url.Parse(t.Announce)
	if err != nil {
		logrus.Errorf("Error parsing tracker url (%v): %v", t.Announce, err)
		return nil, err
	}
	connCtx, connCancel := context.WithCancel(ctx)
	defer connCancel()

	t.UdpManager, err = OpenUdpSocket(connCtx, trackerUrl)
	if err != nil {
		return nil, err
	}

	if err := t.makeConnectUdpReq(); err != nil {
		return nil, err
	}

	//t.makeScrapeUdpReq()

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

	timer := time.NewTimer(time.Second * 4)
	select {
	case <- timer.C:
		return fmt.Errorf("%v tracker udp conn call timed out", t.Announce)
	case data := <- t.UdpManager.Receive:
		body = data
		timer.Stop()
	}

	if body == nil || len(body) < 16 {
		return fmt.Errorf("got invalid body from %v tracker: %v (%v)", t.Announce, body, len(body))
	}

	transId := binary.BigEndian.Uint32(body[4:8])
	if transId != t.TransactionId {
		logrus.Errorf("Tracker resp trans_id (%v) != saved trans_id (%v)", transId, t.TransactionId)
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
		return nil, fmt.Errorf("%v tracker announce call timed out", t.Announce)
	case data := <- t.UdpManager.Receive:
		body = data
		timer.Stop()
	}

	transId := binary.BigEndian.Uint32(body[4:8])
	if transId != t.TransactionId {
		logrus.Errorf("Tracker resp trans id (%v) != saved trans id (%v)", transId, t.TransactionId)
	}
	interval := binary.BigEndian.Uint32(body[8:12])
	leechers := binary.BigEndian.Uint32(body[12:16])
	seeders := binary.BigEndian.Uint32(body[16:20])

	logrus.Infof("Interval = %v; leechers = %v; seeders = %v;", interval, leechers, seeders)
	parsedPeers, err := peers.Unmarshal(body[20:])
	logrus.Infof("Got peers: %v", parsedPeers)
	t.TrackerCallInterval = time.Second * time.Duration(interval)
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
		logrus.Errorf("%v tracker scrape call timed out", t.Announce)
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