package torrentfile

import (
	"encoding/json"
	"net/url"
	"testing"

	"torrentClient/parser/env"

	"github.com/sirupsen/logrus"
)

type testMsg struct {
	PeerPort	uint16	`json:"peer_port"`
	Data		string	`json:"data"`
}

func TestOpenUdpSocket(t *testing.T) {
	tUrl, _ := url.Parse("udp://18.219.47.231:1194")
	conn, err := OpenUdpSocket(tUrl)
	if err != nil {
		t.Errorf("Error opening socket: %v", err)
	}
	testMsg, _ := json.Marshal(testMsg{PeerPort: env.GetParser().GetTorrentPeerPort(), Data: "Hello!"})

	defer func() {
		t.Log("Exiting!")
		conn.ExitChan <- 1
	}()

	conn.Send <- testMsg
	res := <- conn.Receive

	logrus.Infof("Got msg: %v", string(res))

	if string(res) != string("Hello!") {
		t.Errorf("Got not echo msg: %v", string(res))
	}
}

