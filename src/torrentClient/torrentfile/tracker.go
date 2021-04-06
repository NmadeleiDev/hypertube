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

func (t *TorrentFile) requestPeers() ([]peers.Peer, error) {
	if peerIds, err := t.CallFittingScheme(t.Announce); err == nil {
		return peerIds, nil
	} else {
		logrus.Errorf("Error calling main announce: %v", err)
	}

	for _, announce := range t.AnnounceList {
		if peerIds, err := t.CallFittingScheme(announce); err == nil {
			return peerIds, nil
		} else {
			logrus.Errorf("Error calling announce list member: %v", err)
		}
	}

	return nil, fmt.Errorf("failed to call any tracker")
}

func (t *TorrentFile) CallFittingScheme(announce string) ([]peers.Peer, error) {
	trackerUrl, err := url.Parse(announce)
	if err != nil {
		logrus.Errorf("Error parse tracker url: %v", err)
		return nil, err
	}

	if trackerUrl.Scheme == "http" {
		return t.callHttpTracker(announce)
	} else if trackerUrl.Scheme == "udp" {
		return t.callUdpTracker(announce)
	} else {
		return nil, fmt.Errorf("unsupported url scheme: %v; url: %v", trackerUrl.Scheme, t.Announce)
	}
}

func (t *TorrentFile) callUdpTracker(announce string) ([]peers.Peer, error) {
	trackerUrl, err := url.Parse(announce)
	if err != nil {
		logrus.Errorf("Error parsing tracker url (%v): %v", announce, err)
		return nil, err
	}
	t.Download.UdpManager, err = OpenUdpSocket(trackerUrl)
	if err != nil {
		return nil, err
	}

	defer func() {
		t.Download.UdpManager.ExitChan <- 1
	}()

	if err := t.makeConnectUdpReq(); err != nil {
		return nil, err
	}

	t.makeScrapeUdpReq()

	parsedPeers, err := t.makeAnnounceUdpReq()
	return parsedPeers, err
}

func (t *TorrentFile) callHttpTracker(announce string) ([]peers.Peer, error) {
	urlStr, err := t.buildHttpTrackerURL(announce)
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

func (t *TorrentFile) makeConnectUdpReq() error {
	req, err := t.buildUdpTrackerConnectReq()
	if err != nil {
		return err
	}

	t.Download.UdpManager.Send <- req

	var body []byte

	timer := time.NewTimer(time.Second * 3)
	select {
	case <- timer.C:
		return fmt.Errorf("tracker call timed out")
	case data := <- t.Download.UdpManager.Receive:
		body = data
		timer.Stop()
	}

	transId := binary.BigEndian.Uint32(body[4:9])
	if transId != t.Download.TransactionId {
		logrus.Errorf("Tracker resp trans_id (%v) != saved trans_id (%v)", transId, t.Download.TransactionId)
		// выйти?
	}
	t.Download.ConnectionId = binary.BigEndian.Uint64(body[8:])

	logrus.Infof("Connect announce resp: conn_id=%v action=%v trans_id=%v", t.Download.ConnectionId, binary.BigEndian.Uint32(body[:4]), binary.BigEndian.Uint32(body[4:8]))
	return nil
}

func (t *TorrentFile) makeAnnounceUdpReq() ([]peers.Peer, error) {
	req, err := t.buildUdpTrackerAnnounceReq()
	if err != nil {
		return nil, err
	}

	t.Download.UdpManager.Send <- req
	var body []byte

	timer := time.NewTimer(time.Second * 10)

	select {
	case <- timer.C:
		return nil, fmt.Errorf("tracker call timed out")
	case data := <- t.Download.UdpManager.Receive:
		body = data
		timer.Stop()
	}

	transId := binary.BigEndian.Uint32(body[4:8])
	if transId != t.Download.TransactionId {
		logrus.Errorf("Tracker resp trans id (%v) != saved trans id (%v)", transId, t.Download.TransactionId)
		// выйти?
	}
	interval := binary.BigEndian.Uint32(body[8:12])
	leechers := binary.BigEndian.Uint32(body[12:16])
	seeders := binary.BigEndian.Uint32(body[16:20])

	logrus.Infof("Interval = %v; leechers = %v; seeders = %v;", interval, leechers, seeders)
	parsedPeers, err := peers.Unmarshal(body[20:])
	logrus.Infof("Got peers: %v", parsedPeers)
	return parsedPeers, err
}

func (t *TorrentFile) makeScrapeUdpReq() {
	req, err := t.buildScrapeUdpReq()
	if err != nil {
		logrus.Errorf("Error building scrape req: %v", err)
		return
	}

	t.Download.UdpManager.Send <- req
	var body []byte

	timer := time.NewTimer(time.Second * 10)

	select {
	case <- timer.C:
		logrus.Errorf("tracker call timed out")
		return
	case data := <- t.Download.UdpManager.Receive:
		body = data
		timer.Stop()
	}

	transId := binary.BigEndian.Uint32(body[4:8])
	if transId != t.Download.TransactionId {
		logrus.Errorf("Tracker resp trans id (%v) != saved trans id (%v)", transId, t.Download.TransactionId)
		// выйти?
	}
	seeders := binary.BigEndian.Uint32(body[8:12])
	completed := binary.BigEndian.Uint32(body[12:16])
	leechers := binary.BigEndian.Uint32(body[16:20])

	logrus.Infof("Scrape res: completed = %v; leechers = %v; seeders = %v;", completed, leechers, seeders)
	return
}