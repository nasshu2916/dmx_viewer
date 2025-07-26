package artnet

import (
	"errors"
	"time"
)

// 定数定義
const (
	DefaultPort              = 6454
	DefaultChannelBufferSize = 1000
	DefaultReadTimeout       = 500 * time.Millisecond
	DefaultStatInterval      = 60 * time.Second
	DefaultMaxPacketSize     = 1500

	// チャンネル使用率の閾値
	HighUtilizationThreshold     = 75.0
	CriticalUtilizationThreshold = 90.0
)

// エラー定義
var (
	ErrSendChannelFull          = errors.New("send channel is full")
	ErrSendChannelTimeout       = errors.New("send channel write timeout")
	ErrConnectionNotEstablished = errors.New("connection is not established")
	ErrServerNotRunning         = errors.New("server is not running")
)

// ChannelType チャンネルタイプ
type ChannelType int

const (
	ReceiveChannel ChannelType = iota
	SendChannel
)

func (ct ChannelType) String() string {
	switch ct {
	case ReceiveChannel:
		return "receive"
	case SendChannel:
		return "send"
	default:
		return "unknown"
	}
}

// HealthStatus 健全性ステータス
type HealthStatus int

const (
	HealthyStatus HealthStatus = iota
	WarningStatus
	CriticalStatus
)

func (hs HealthStatus) String() string {
	switch hs {
	case HealthyStatus:
		return "healthy"
	case WarningStatus:
		return "warning"
	case CriticalStatus:
		return "critical"
	default:
		return "unknown"
	}
}
