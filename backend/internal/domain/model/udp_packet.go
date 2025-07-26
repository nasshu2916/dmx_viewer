package model

import (
	"net"

	"github.com/jsimonetti/go-artnet/packet"
)

type ReceivedData struct {
	Data []byte
	Addr net.Addr
}

type ReceivedArtPacket struct {
	Packet packet.ArtNetPacket
	Addr   net.Addr
}
