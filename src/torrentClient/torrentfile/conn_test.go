package torrentfile

import (
	"bytes"
	"net/url"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestOpenUdpSocket(t *testing.T) {
	tUrl, _ := url.Parse("udp://18.219.47.231:1194")
	conn, err := OpenUdpSocket(tUrl)
	if err != nil {
		t.Errorf("Error opening socket: %v", err)
	}
	testMsg := []byte("Hello motherfucker!")
	defer func() {
		t.Log("Exiting!")
		conn.ExitChan <- 1
	}()

	resp := bytes.ReplaceAll(testMsg, []byte("fuck"), []byte("shit"))

	conn.Send <- testMsg
	res := <- conn.Receive

	logrus.Infof("Got msg: %v", string(res))

	if string(res) != string(resp) {
		t.Errorf("Got not echo msg: %v", string(res))
	}
}

