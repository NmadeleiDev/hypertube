package client

import (
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
