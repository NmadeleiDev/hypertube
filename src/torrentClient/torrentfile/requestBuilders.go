package torrentfile

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"

	"github.com/sirupsen/logrus"
)

const (
	defaultAnnouncePeers = 50
)

func (t *Tracker) buildHttpTrackerURL(announce string) (string, error) {
	base, err := url.Parse(announce)
	if err != nil {
		return "", fmt.Errorf("url parse error: %s; src: %s", err.Error(), t.Announce)
	}
	params := url.Values{
		"info_hash":  []string{string(t.InfoHash[:])},
		"peer_id":    []string{string(t.MyPeerId[:])},
		"port":       []string{strconv.Itoa(int(t.MyPeerPort))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(t.Length)},
	}
	base.RawQuery = params.Encode()
	return base.String(), nil
}

func (t *Tracker) buildUdpTrackerConnectReq() (request []byte, err error) {
	t.TransactionId = uint32(rand.Int31())

	req := make([]byte, 16)
	binary.BigEndian.PutUint64(req[:8], uint64(protocolId))
	binary.BigEndian.PutUint32(req[8:12], uint32(connectAction))
	binary.BigEndian.PutUint32(req[12:], t.TransactionId)

	//logrus.Infof("Ready msg: %v; %v, %v, %v, %v", req, string(req), fmt.Sprint(req[:8]), fmt.Sprint(req[8:12]), fmt.Sprint(req[12:]))
	logrus.Infof("Conn to server msg: %v %v %v", fmt.Sprint(binary.BigEndian.Uint64(req[:8])), fmt.Sprint(binary.BigEndian.Uint32(req[8:12])), fmt.Sprint(binary.BigEndian.Uint32(req[12:])))
	return req, nil
}

func (t *Tracker) buildUdpTrackerAnnounceReq() (request []byte, err error) {
	t.TransactionId = uint32(rand.Int31())

	req := make([]byte, 98)

	//fmt.Printf("Infohash:\n(%b) \n(%s) \n(%v)\n", t.InfoHash[:], string(t.InfoHash[:]), t.InfoHash[:])

	binary.BigEndian.PutUint64(req[0:8], t.ConnectionId)
	binary.BigEndian.PutUint32(req[8:12], uint32(announceAction))
	binary.BigEndian.PutUint32(req[12:16], t.TransactionId)
	copy(req[16:36], t.InfoHash[:])
	copy(req[36:56], t.MyPeerId[:])
	binary.BigEndian.PutUint64(req[56:64], 0) // downloaded
	binary.BigEndian.PutUint64(req[64:72], uint64(t.PieceLength * len(t.PieceHashes)))
	binary.BigEndian.PutUint64(req[72:80], 0)                    // uploaded
	binary.BigEndian.PutUint32(req[80:84], 0)                    // event
	binary.BigEndian.PutUint32(req[84:88], 0)                    // ip addr
	binary.BigEndian.PutUint32(req[88:92], 0) // key
	binary.BigEndian.PutUint32(req[92:96], uint32(defaultAnnouncePeers)) // num want
	binary.BigEndian.PutUint16(req[96:98], t.MyPeerPort)

	logrus.Infof("my peer id: %v; left: %v; my port: %v", t.MyPeerId, uint64(t.PieceLength * len(t.PieceHashes)), t.MyPeerPort)

	return req, nil
}

func (t *Tracker) buildScrapeUdpReq() (request []byte, err error) {
	t.TransactionId = uint32(rand.Int31())

	req := make([]byte, 16 + len(t.InfoHash))
	binary.BigEndian.PutUint64(req[:8], t.ConnectionId)
	binary.BigEndian.PutUint32(req[8:12], uint32(scrapeAction))
	binary.BigEndian.PutUint32(req[12:16], t.TransactionId)
	copy(req[16:], t.InfoHash[:])

	return req, nil
}
