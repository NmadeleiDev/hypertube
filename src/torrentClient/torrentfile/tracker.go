package torrentfile

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"torrent_client/peers"

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
)

func (t *TorrentFile) buildHttpTrackerURL(peerID [20]byte, port uint16) (string, error) {
	base, err := url.Parse(t.Announce)
	if err != nil {
		return "", fmt.Errorf("url parse error: %s; src: %s", err.Error(), t.Announce)
	}
	params := url.Values{
		"info_hash":  []string{string(t.InfoHash[:])},
		"peer_id":    []string{string(peerID[:])},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(t.Length)},
	}
	base.RawQuery = params.Encode()
	return base.String(), nil
}

func (t *TorrentFile) buildUdpTrackerConnectReq() (request []byte, err error) {
	t.Download.TransactionId = uint32(rand.Int31())

	req := make([]byte, 16)
	binary.BigEndian.PutUint64(req[:9], uint64(protocolId))
	binary.BigEndian.PutUint32(req[8:13], uint32(connectAction))
	binary.BigEndian.PutUint32(req[12:], t.Download.TransactionId)

	return req, nil
}

func (t *TorrentFile) buildUdpTrackerAnnounceReq(peerID [20]byte, port uint16) (request []byte, err error) {
	t.Download.TransactionId = uint32(rand.Int31())

	req := make([]byte, 16)

	binary.BigEndian.PutUint64(req[:9], uint64(t.Download.ConnectionId))
	binary.BigEndian.PutUint32(req[8:13], uint32(announceAction))
	binary.BigEndian.PutUint32(req[12:17], t.Download.TransactionId)
	copy(req[16:37], t.InfoHash[:])
	copy(req[36:57], peerID[:])
	binary.BigEndian.PutUint64(req[56:65], 0)
	binary.BigEndian.PutUint64(req[64:73], uint64(t.Length))
	binary.BigEndian.PutUint64(req[72:81], 0) // uploaded
	binary.BigEndian.PutUint32(req[80:85], 0) // event
	binary.BigEndian.PutUint32(req[84:89], 0) // ip addr
	binary.BigEndian.PutUint32(req[88:93], 0) // key
	binary.PutVarint(req[92:97], -1) // num want
	binary.BigEndian.PutUint16(req[96:99], port)

	return req, nil
}

func (t *TorrentFile) requestPeers(peerID [20]byte, port uint16) ([]peers.Peer, error) {
	trackerUrl, err := url.Parse(t.Announce)
	if err != nil {
		logrus.Errorf("Error parse tracker url: %v", err)
		return nil, err
	}

	if trackerUrl.Scheme == "http" {
		return t.callHttpTracker(peerID, port)
	} else if trackerUrl.Scheme == "udp" {
		return t.callUdpTracker(peerID, port)
	} else {
		return nil, fmt.Errorf("unsupported url scheme: %v; url: %v", trackerUrl.Scheme, t.Announce)
	}
}

func (t *TorrentFile) callUdpTracker(peerID [20]byte, port uint16) ([]peers.Peer, error) {
	trackerUrl, err := url.Parse(t.Announce)
	if err != nil {
		logrus.Errorf("Error parsing tracker url (%v): %v", t.Announce, err)
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

	return t.makeAnnounceUdpReq(peerID, port)
}

func (t *TorrentFile) callHttpTracker(peerID [20]byte, port uint16) ([]peers.Peer, error) {
	urlStr, err := t.buildHttpTrackerURL(peerID, port)
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

	timer := time.NewTimer(time.Second * 10)
	select {
	case <- timer.C:
		return fmt.Errorf("tracker call timed out")
	case data := <- t.Download.UdpManager.Receive:
		body = data
	}

	transId := binary.BigEndian.Uint32(body[4:9])
	if transId != t.Download.TransactionId {
		logrus.Errorf("Tracker resp trans id (%v) != saved trans id (%v)", transId, t.Download.TransactionId)
		// выйти?
	}
	t.Download.ConnectionId = binary.BigEndian.Uint64(body[8:17])
	return nil
}

func (t *TorrentFile) makeAnnounceUdpReq(peerID [20]byte, port uint16) ([]peers.Peer, error) {
	req, err := t.buildUdpTrackerAnnounceReq(peerID, port)
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
	}

	transId := binary.BigEndian.Uint32(body[4:9])
	if transId != t.Download.TransactionId {
		logrus.Errorf("Tracker resp trans id (%v) != saved trans id (%v)", transId, t.Download.TransactionId)
		// выйти?
	}
	interval := binary.BigEndian.Uint32(body[8:13])
	leechers := binary.BigEndian.Uint32(body[12:17])
	seeders := binary.BigEndian.Uint32(body[16:21])

	logrus.Infof("Interval = %v; leechers = %v; seeders = %v;", interval, leechers, seeders)
	return peers.Unmarshal(body[20:])
}