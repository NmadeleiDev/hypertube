package client

import (
	"fmt"
	"net"

	"torrentClient/bitfield"
	"torrentClient/peers"
)

type Client struct {
	Conn     net.Conn
	Choked   bool
	Bitfield bitfield.Bitfield
	peer     peers.Peer
	infoHash [20]byte
	peerID   [20]byte
}

func (c *Client) GetClientInfo() string {
	return fmt.Sprintf("Peer: %v\nChoked = %v\nBitfield: %v\n", c.peer.GetAddr(), c.Conn, c.Bitfield)
}
