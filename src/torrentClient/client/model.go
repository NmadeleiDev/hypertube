package client

import (
	"fmt"
	"net"
	"sync"

	"torrentClient/bitfield"
	"torrentClient/peers"
)

type Client struct {
	Mu       sync.Mutex
	Conn     net.Conn
	Choked   bool
	Bitfield bitfield.Bitfield
	Peer     peers.Peer
	infoHash [20]byte
	peerID   [20]byte
}

func (c *Client) GetClientInfo() string {
	return fmt.Sprintf("Peer: %v\nChoked = %v\nBitfield: %v\n", c.peer.GetAddr(), c.Choked, c.Bitfield)
}

func (c *Client) GetShortInfo() string {
	return fmt.Sprintf("Peer addr: %v, is choked = %v", c.peer.GetAddr(), c.Choked)
}

func (c *Client) GetPeer() peers.Peer {
	return c.peer
}
