package model

import (
	"errors"
	"fmt"
	"net"

	"github.com/jsimonetti/go-artnet/packet"
)

const (
	MaxUniverse = 0x7FFF
)

// DMXData DMXデータを表すドメインモデル
type DMXData struct {
	Sequence   uint8      `json:"Sequence"`   // シーケンス番号
	Physical   uint8      `json:"Physical"`   // 物理出力ポート
	SubUni     uint8      `json:"SubUni"`     // サブユニバース（下位8ビット）
	Net        uint8      `json:"Net"`        // ネット（上位8ビット）
	Length     uint16     `json:"Length"`     // データ長
	Data       [512]uint8 `json:"Data"`       // DMXチャンネルデータ
	SourceIP   net.IP     `json:"SourceIP"`   // 送信元IPアドレス
	SourcePort int        `json:"SourcePort"` // 送信元ポート番号
}

// NewDMXData ArtDMXPacketからDMXDataを作成
func NewDMXData(srcAddr net.Addr, packet *packet.ArtDMXPacket) (*DMXData, error) {
	if packet == nil {
		return nil, errors.New("packet cannot be nil")
	}

	addr := srcAddr.(*net.UDPAddr)

	dmx := &DMXData{
		Sequence:   packet.Sequence,
		Physical:   packet.Physical,
		SubUni:     packet.SubUni,
		Net:        packet.Net,
		Length:     packet.Length,
		Data:       packet.Data,
		SourceIP:   addr.IP,
		SourcePort: addr.Port,
	}

	if err := dmx.Validate(); err != nil {
		return nil, fmt.Errorf("invalid DMX data: %w", err)
	}

	return dmx, nil
}

// Validate DMXDataの妥当性を検証
func (d *DMXData) Validate() error {
	if d.Length > 512 {
		return fmt.Errorf("length %d exceeds maximum DMX channels %d", d.Length, 512)
	}

	universe := d.GetUniverse()
	if universe > MaxUniverse {
		return fmt.Errorf("universe %d exceeds maximum %d", universe, MaxUniverse)
	}

	return nil
}

// GetUniverse ネットとサブユニバースからユニバース番号を計算
func (d *DMXData) GetUniverse() uint16 {
	return (uint16(d.Net) << 8) | uint16(d.SubUni)
}

// SetUniverse ユニバース番号からネットとサブユニバースを設定
func (d *DMXData) SetUniverse(universe uint16) {
	d.Net = uint8(universe >> 8)
	d.SubUni = uint8(universe & 0xFF)
}

// 指定チャンネルの値を取得（1-based）
func (d *DMXData) GetChannelValue(channel int) (uint8, error) {
	if channel < 1 || channel > 512 {
		return 0, fmt.Errorf("channel %d out of range (1-512)", channel)
	}

	index := channel - 1
	if uint16(index) >= d.Length {
		return 0, nil // 範囲外は0を返す
	}

	return d.Data[index], nil
}

// 指定チャンネルの値を設定（1-based）
func (d *DMXData) SetChannelValue(channel int, value uint8) error {
	if channel < 1 || channel > 512 {
		return fmt.Errorf("channel %d out of range (1-512)", channel)
	}

	index := channel - 1
	d.Data[index] = value

	// Lengthを必要に応じて拡張
	if uint16(channel) > d.Length {
		d.Length = uint16(channel)
	}

	return nil
}

// 指定範囲のチャンネル値を取得（1-based）
func (d *DMXData) GetChannelRange(startChannel, endChannel int) ([]uint8, error) {
	if startChannel < 1 || startChannel > 512 {
		return nil, fmt.Errorf("start channel %d out of range (1-512)", startChannel)
	}
	if endChannel < 1 || endChannel > 512 {
		return nil, fmt.Errorf("end channel %d out of range (1-512)", endChannel)
	}
	if startChannel > endChannel {
		return nil, fmt.Errorf("start channel %d cannot be greater than end channel %d", startChannel, endChannel)
	}

	result := make([]uint8, endChannel-startChannel+1)
	for i := startChannel; i <= endChannel; i++ {
		value, _ := d.GetChannelValue(i)
		result[i-startChannel] = value
	}

	return result, nil
}

// String DMXDataの文字列表現
func (d *DMXData) String() string {
	return fmt.Sprintf("DMX[Universe:%d, Seq:%d, Length:%d",
		d.GetUniverse(), d.Sequence, d.Length)
}

// DMXDataの深いコピーを作成
func (d *DMXData) Clone() *DMXData {
	clone := &DMXData{
		Sequence:   d.Sequence,
		Physical:   d.Physical,
		SubUni:     d.SubUni,
		Net:        d.Net,
		Length:     d.Length,
		SourceIP:   make(net.IP, len(d.SourceIP)),
		SourcePort: d.SourcePort,
	}
	copy(clone.Data[:], d.Data[:])
	return clone
}
