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

func (t *TorrentFile) buildHttpTrackerURL(announce string, peerID [20]byte, port uint16) (string, error) {
	base, err := url.Parse(announce)
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

	//logrus.Infof("Connecting with: [%v, %v, %v]", protocolId, connectAction, t.Download.TransactionId)
	req := make([]byte, 16)
	binary.BigEndian.PutUint64(req[:8], uint64(protocolId))
	binary.BigEndian.PutUint32(req[8:12], uint32(connectAction))
	binary.BigEndian.PutUint32(req[12:], t.Download.TransactionId)

	logrus.Infof("Ready msg: %v; %v, %v, %v, %v", req, string(req), fmt.Sprint(req[:8]), fmt.Sprint(req[8:12]), fmt.Sprint(req[12:]))
	logrus.Infof("Str form: %v %v %v", fmt.Sprint(binary.BigEndian.Uint64(req[:8])), fmt.Sprint(binary.BigEndian.Uint32(req[8:12])), fmt.Sprint(binary.BigEndian.Uint32(req[12:])))
	return req, nil
}

func (t *TorrentFile) buildUdpTrackerAnnounceReq(peerID [20]byte, port uint16) (request []byte, err error) {
	t.Download.TransactionId = uint32(rand.Int31())

	req := make([]byte, 98)

	binary.BigEndian.PutUint64(req[:8], uint64(t.Download.ConnectionId))
	binary.BigEndian.PutUint32(req[8:12], uint32(announceAction))
	binary.BigEndian.PutUint32(req[12:16], t.Download.TransactionId)
	copy(req[16:36], t.InfoHash[:])
	copy(req[36:56], peerID[:])
	binary.BigEndian.PutUint64(req[56:64], 0)
	binary.BigEndian.PutUint64(req[64:72], uint64(t.Length))
	binary.BigEndian.PutUint64(req[72:80], 0) // uploaded
	binary.BigEndian.PutUint32(req[80:84], 0) // event
	binary.BigEndian.PutUint32(req[84:88], 0) // ip addr
	binary.BigEndian.PutUint32(req[88:92], 0) // key
	binary.PutVarint(req[92:96], -1) // num want
	binary.BigEndian.PutUint16(req[96:98], port)

	return req, nil
}

func (t *TorrentFile) requestPeers(peerID [20]byte, port uint16) ([]peers.Peer, error) {
	if peerIds, err := t.CallFittingScheme(t.Announce, peerID, port); err == nil {
		return peerIds, nil
	} else {
		logrus.Errorf("Error calling main announce: %v", err)
	}

	for _, announce := range t.AnnounceList {
		if peerIds, err := t.CallFittingScheme(announce, peerID, port); err == nil {
			return peerIds, nil
		} else {
			logrus.Errorf("Error calling announce list member: %v", err)
		}
	}

	return nil, fmt.Errorf("failed to call any tracker")
}

func (t *TorrentFile) CallFittingScheme(announce string, peerID [20]byte, port uint16) ([]peers.Peer, error) {
	trackerUrl, err := url.Parse(announce)
	if err != nil {
		logrus.Errorf("Error parse tracker url: %v", err)
		return nil, err
	}

	if trackerUrl.Scheme == "http" {
		return t.callHttpTracker(announce, peerID, port)
	} else if trackerUrl.Scheme == "udp" {
		return t.callUdpTracker(announce, peerID, port)
	} else {
		return nil, fmt.Errorf("unsupported url scheme: %v; url: %v", trackerUrl.Scheme, t.Announce)
	}
}

func (t *TorrentFile) callUdpTracker(announce string, peerID [20]byte, port uint16) ([]peers.Peer, error) {
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
		//logrus.Info("Writing exit msg")
		t.Download.UdpManager.ExitChan <- 1
	}()

	if err := t.makeConnectUdpReq(); err != nil {
		return nil, err
	}

	parsedPeers, err := t.makeAnnounceUdpReq(peerID, port)
	return parsedPeers, err
}

func (t *TorrentFile) callHttpTracker(announce string, peerID [20]byte, port uint16) ([]peers.Peer, error) {
	urlStr, err := t.buildHttpTrackerURL(announce, peerID, port)
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

	logrus.Infof("Writing %v bytes as conn req", len(req))
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

func (t *TorrentFile) makeAnnounceUdpReq(peerID [20]byte, port uint16) ([]peers.Peer, error) {
	req, err := t.buildUdpTrackerAnnounceReq(peerID, port)
	if err != nil {
		return nil, err
	}

	logrus.Infof("Writing %v bytes as conn announce", len(req))
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