package torrentfile

import (
	"net"
	"net/url"

	"github.com/sirupsen/logrus"
)


func StartHandlingSocket(conn *net.UDPConn, udpManager *UdpConnManager)  {
	go func() {
		for {
			if !udpManager.IsValid() {
				return
			}

			buffer := make([]byte, 1024)
			n, _, err := conn.ReadFromUDP(buffer)
			if err != nil {
				logrus.Errorf("Error reading from conn: %v", err)
				return
			} else {
				logrus.Infof("Read %v bytes", n)
				udpManager.Receive <- buffer[:n]
			}
		}
	}()

	go func() {
		for {
			msg := <- udpManager.Send
			if !udpManager.IsValid() {
				return
			}
			if n, err := conn.Write(msg); err != nil {
				logrus.Errorf("Error write msg: %v", err)
				return
			} else {
				logrus.Infof("Wrote %v bytes", n)
			}
		}
	}()
}

func OpenUdpSocket(tUrl *url.URL) (*UdpConnManager, error) {
	addr := tUrl.Host
	destinationAddress, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		logrus.Errorf("Error resolving addr %v : %v", addr, err)
		return nil, err
	}

	connection, err := net.DialUDP("udp", nil, destinationAddress)
	if err != nil {
		logrus.Errorf("Error dial addr %v : %v",  err, addr)
		return nil, err
	} else {
		logrus.Infof("Connection with %v opened (%v)", *destinationAddress, tUrl.String())
	}

	exitChan := make(chan byte)
	receive := make(chan []byte, 10)
	send := make(chan []byte, 10)

	manager := &UdpConnManager{Receive: receive, Send: send, ExitChan: exitChan, isValid: true}

	go func() {
		<- exitChan
		logrus.Info("Exiting from udp socket.")
		manager.SetValid(false)
		connection.Close()
		close(receive)
		close(send)
		close(exitChan)
	}()

	StartHandlingSocket(connection, manager)
	return manager, nil
}
