package model

import (
	"net"
	"testing"

	"github.com/jsimonetti/go-artnet/packet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDMXData(t *testing.T) {
	tests := []struct {
		name    string
		packet  *packet.ArtDMXPacket
		wantErr bool
	}{
		{
			name: "Valid DMX packet",
			packet: &packet.ArtDMXPacket{
				Sequence: 1,
				Physical: 0,
				SubUni:   5,
				Net:      2,
				Length:   100,
				Data:     [512]byte{255, 128, 64, 32},
			},
			wantErr: false,
		},
		{
			name:    "Nil packet",
			packet:  nil,
			wantErr: true,
		},
		{
			name: "Invalid length",
			packet: &packet.ArtDMXPacket{
				Length: 513, // 最大値を超える
			},
			wantErr: true,
		},
	}

	addr := &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 1234,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dmx, err := NewDMXData(addr, tt.packet)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, dmx)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, dmx)
				assert.Equal(t, tt.packet.Sequence, dmx.Sequence)
				assert.Equal(t, tt.packet.Physical, dmx.Physical)
				assert.Equal(t, tt.packet.SubUni, dmx.SubUni)
				assert.Equal(t, tt.packet.Net, dmx.Net)
				assert.Equal(t, tt.packet.Length, dmx.Length)
				assert.Equal(t, tt.packet.Data, dmx.Data)
			}
		})
	}
}

func TestDMXData_GetUniverse(t *testing.T) {
	dmx := &DMXData{
		Net:    2,
		SubUni: 5,
	}

	expected := uint16(2<<8) | uint16(5) // 2*256 + 5 = 517
	assert.Equal(t, expected, dmx.GetUniverse())
}

func TestDMXData_SetUniverse(t *testing.T) {
	dmx := &DMXData{}
	universe := uint16(517) // 2*256 + 5

	dmx.SetUniverse(universe)

	assert.Equal(t, uint8(2), dmx.Net)
	assert.Equal(t, uint8(5), dmx.SubUni)
}

func TestDMXData_GetChannelValue(t *testing.T) {
	dmx := &DMXData{
		Length: 10,
		Data:   [512]byte{255, 128, 64, 32, 16, 8, 4, 2, 1, 0},
	}

	tests := []struct {
		name     string
		channel  int
		expected uint8
		wantErr  bool
	}{
		{
			name:     "Valid channel 1",
			channel:  1,
			expected: 255,
			wantErr:  false,
		},
		{
			name:     "Valid channel 5",
			channel:  5,
			expected: 16,
			wantErr:  false,
		},
		{
			name:     "Channel beyond length",
			channel:  15,
			expected: 0, // 範囲外は0
			wantErr:  false,
		},
		{
			name:    "Invalid channel 0",
			channel: 0,
			wantErr: true,
		},
		{
			name:    "Invalid channel 513",
			channel: 513,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := dmx.GetChannelValue(tt.channel)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, value)
			}
		})
	}
}

func TestDMXData_SetChannelValue(t *testing.T) {
	dmx := &DMXData{
		Length: 5,
	}

	// 有効なチャンネルに値を設定
	err := dmx.SetChannelValue(3, 128)
	assert.NoError(t, err)
	assert.Equal(t, uint8(128), dmx.Data[2]) // 3番目のチャンネルは配列のインデックス2

	// Lengthを超えるチャンネルに設定（Lengthが拡張される）
	err = dmx.SetChannelValue(10, 64)
	assert.NoError(t, err)
	assert.Equal(t, uint8(64), dmx.Data[9])
	assert.Equal(t, uint16(10), dmx.Length)

	// 無効なチャンネル
	err = dmx.SetChannelValue(0, 255)
	assert.Error(t, err)

	err = dmx.SetChannelValue(513, 255)
	assert.Error(t, err)
}

func TestDMXData_GetChannelRange(t *testing.T) {
	dmx := &DMXData{
		Length: 10,
		Data:   [512]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}

	// 有効な範囲
	values, err := dmx.GetChannelRange(3, 6)
	require.NoError(t, err)
	expected := []uint8{3, 4, 5, 6}
	assert.Equal(t, expected, values)

	// 無効な範囲
	_, err = dmx.GetChannelRange(0, 5)
	assert.Error(t, err)

	_, err = dmx.GetChannelRange(5, 3) // start > end
	assert.Error(t, err)
}

func TestDMXData_Clone(t *testing.T) {
	original := &DMXData{
		Sequence: 1,
		Physical: 2,
		SubUni:   3,
		Net:      4,
		Length:   10,
		Data:     [512]byte{1, 2, 3, 4, 5},
	}

	clone := original.Clone()

	// 値が同じことを確認
	assert.Equal(t, original.Sequence, clone.Sequence)
	assert.Equal(t, original.Physical, clone.Physical)
	assert.Equal(t, original.SubUni, clone.SubUni)
	assert.Equal(t, original.Net, clone.Net)
	assert.Equal(t, original.Length, clone.Length)
	assert.Equal(t, original.Data, clone.Data)

	// 異なるオブジェクトであることを確認
	assert.NotSame(t, original, clone)

	// 一方を変更しても他方に影響しないことを確認
	clone.Sequence = 99
	clone.Data[0] = 99
	assert.NotEqual(t, original.Sequence, clone.Sequence)
	assert.NotEqual(t, original.Data[0], clone.Data[0])
}

func TestDMXData_String(t *testing.T) {
	dmx := &DMXData{
		Net:      2,
		SubUni:   5,
		Sequence: 10,
		Length:   100,
		Data:     [512]byte{255, 128, 64}, // 3つのアクティブチャンネル
	}

	str := dmx.String()
	expected := "DMX[Universe:517, Seq:10, Length:100]" // 2*256+5 = 517
	assert.Equal(t, expected, str)
}
