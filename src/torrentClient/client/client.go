package client

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"time"

	"torrentClient/bitfield"
	"torrentClient/handshake"
	"torrentClient/message"
	"torrentClient/peers"

	"github.com/sirupsen/logrus"
)


func completeHandshake(conn net.Conn, infohash, peerID [20]byte) (*handshake.Handshake, error) {
	conn.SetDeadline(time.Now().Add(3 * time.Second))
	defer conn.SetDeadline(time.Time{}) // Disable the deadline

	req := handshake.New(infohash, peerID)
	_, err := conn.Write(req.Serialize())
	if err != nil {
		return nil, fmt.Errorf("request write error: %v", err)
	} else {
		logrus.Debugf("Wrote handshake msg (%v bytes)", len(req.Serialize()))
	}

	res, err := handshake.Read(conn)
	if err != nil {
		return nil, fmt.Errorf("read error: %v", err)
	}
	if !bytes.Equal(res.InfoHash[:], infohash[:]) {
		return nil, fmt.Errorf("expected infohash %x but got %x", res.InfoHash, infohash)
	}
	return res, nil
}

func recvBitfield(conn net.Conn) (bitfield.Bitfield, error) {
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	defer conn.SetDeadline(time.Time{}) // Disable the deadline

	msg, err := message.Read(conn)
	if err != nil {
		return nil, err
	}
	if msg == nil {
		err := fmt.Errorf("expected bitfield but got %s", msg)
		return nil, err
	}
	if msg.ID != message.MsgBitfield {
		err := fmt.Errorf("expected bitfield but got ID %d", msg.ID)
		return nil, err
	}

	return msg.Payload, nil
}

func New(peer peers.Peer, peerID, infoHash [20]byte) (*Client, error) {
	conn, err := net.DialTimeout("tcp", peer.GetAddr(), 10 * time.Second)
	if err != nil {
		return nil, fmt.Errorf("dial error: %v; was connecting to %v", err, peer.GetAddr())
	} else {
		logrus.Debugf("Connected to Peer on %v", peer.GetAddr())
	}

	_, err = completeHandshake(conn, infoHash, peerID)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("handshake error: %v", err)
	}

	bf, err := recvBitfield(conn)
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Client{
		Conn:     conn,
		Choked:   true,
		Bitfield: bf,
		Peer:     peer,
		infoHash: infoHash,
		peerID:   peerID,
	}, nil
}

func (c *Client) Read() (*message.Message, error) {
	msg, err := message.Read(c.Conn)
	return msg, err
}

func (c *Client) WaitForUnchoke(ctx context.Context) (bool, error) {
	defer c.Conn.SetDeadline(time.Time{}) // Disable the deadline
	tries := 0

	for {
		select {
		case <- ctx.Done():
			return false, fmt.Errorf("exited because of ctx done")
		default:
			if tries > 10 {
				return false, fmt.Errorf("too many tries (%v) to init peer %v", tries, c.Peer.GetAddr())
			}
			if err := c.Conn.SetDeadline(time.Now().Add(15 * time.Second)); err != nil {
				logrus.Debugf("Error setting deadline for conn for init peer %v", c.Peer.GetAddr())
				c.Peer.IsDead = true
				tries++
				return false, err
			}

			if err := c.SendUnchoke(); err != nil {
				logrus.Errorf("Error sending unchoke: %v", err)
				c.Peer.IsDead = true
				tries++
				return false, fmt.Errorf("failed to send unchoke: %v", err)
			}

			msg, err := c.Read()
			if err != nil {
				logrus.Debugf("Error waiting for peer %v msg in unchoke wait: %v", c.Peer.GetAddr(), err)
				tries++
				continue
				//return false, err
			}

			if msg == nil { // keep alive
				continue
			}

			switch msg.ID {
			case message.MsgUnchoke:
				c.Choked = false
				logrus.Infof("Got UNCHOKE from %v", c.GetShortInfo())
				return true, nil
			case message.MsgChoke:
				c.Choked = true
			case message.MsgHave:
				index, err := message.ParseHave(msg)
				if err != nil {
					return false, err
				}
				c.Bitfield.SetPiece(index)
			default:
				logrus.Debugf("Got unexpeced msg id=%v from peer %v in init", msg.ID, c.Peer.GetAddr())
			}
		}
	}
}

func (c *Client) SendRequest(index, begin, length int) error {
	req := message.FormatRequest(index, begin, length)
	_, err := c.Conn.Write(req.Serialize())
	return err
}

func (c *Client) SendInterested() error {
	msg := message.Message{ID: message.MsgInterested}
	_, err := c.Conn.Write(msg.Serialize())
	if err != nil {
		logrus.Errorf("Error sending interested msg: %v", err)
	}
	return err
}

func (c *Client) SendKeepAlive() error {
	if _, err := c.Conn.Write(make([]byte, 4)); err != nil {
		logrus.Errorf("Error sending keep alive: %v", err)
		return err
	}
	return nil
}

func (c *Client) SendNotInterested() error {
	msg := message.Message{ID: message.MsgNotInterested}
	_, err := c.Conn.Write(msg.Serialize())
	if err != nil {
		logrus.Errorf("Error sending not interested msg: %v", err)
	}
	return err
}

func (c *Client) SendUnchoke() error {
	msg := message.Message{ID: message.MsgUnchoke}
	_, err := c.Conn.Write(msg.Serialize())
	if err != nil {
		logrus.Errorf("Error sending unchoke msg: %v", err)
	}
	return err
}

func (c *Client) SendHave(index int) error {
	msg := message.FormatHave(index)
	_, err := c.Conn.Write(msg.Serialize())
	if err != nil {
		logrus.Errorf("Error sending have msg: %v", err)
	}
	return err
}
