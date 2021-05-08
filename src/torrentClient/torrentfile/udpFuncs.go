package torrentfile

import (
	"context"
	"net"
	"net/url"

	"github.com/sirupsen/logrus"
)


func StartHandlingSocket(ctx context.Context, conn *net.UDPConn, udpManager *UdpConnManager)  {
	go func() {
		defer close(udpManager.Receive)

		for {
			select {
			case <- ctx.Done():
				return
			default:
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
		}
	}()

	go func() {
		defer close(udpManager.Send)

		for {
			select {
			case <- ctx.Done():
				return
			case msg := <- udpManager.Send:
				if n, err := conn.Write(msg); err != nil {
					logrus.Errorf("Error write msg: %v", err)
					return
				} else {
					logrus.Infof("Wrote %v bytes", n)
				}
			}
		}
	}()
}

func OpenUdpSocket(ctx context.Context, tUrl *url.URL) (*UdpConnManager, error) {
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

	receive := make(chan []byte, 10)
	send := make(chan []byte, 10)

	manager := &UdpConnManager{Receive: receive, Send: send}

	StartHandlingSocket(ctx, connection, manager)
	return manager, nil
}
