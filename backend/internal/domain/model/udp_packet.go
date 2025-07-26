package model

import "net"

type ReceivedPacket struct {
	Data []byte
	Addr net.Addr
}
