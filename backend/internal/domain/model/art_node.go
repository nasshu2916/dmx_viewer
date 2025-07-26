package model

import (
	"bytes"
	"net"
	"time"

	"github.com/jsimonetti/go-artnet/packet"
)

type ArtNetNode struct {
	IPAddress  net.IP
	ShortName  string
	LongName   string
	NodeReport string
	MacAddress net.HardwareAddr
	LastSeen   time.Time
}

func NewArtNetNode(p *packet.ArtPollReplyPacket) *ArtNetNode {
	shortName := string(bytes.Trim(p.ShortName[:], "\x00"))
	longName := string(bytes.Trim(p.LongName[:], "\x00"))
	nodeReport := string(bytes.Trim([]byte(string(p.NodeReport[:])), "\x00"))

	return &ArtNetNode{
		IPAddress:  net.IP(p.IPAddress[:]),
		ShortName:  shortName,
		LongName:   longName,
		NodeReport: nodeReport,
		MacAddress: net.HardwareAddr(p.Macaddress[:]),
		LastSeen:   time.Now(),
	}
}
