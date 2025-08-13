package artnet

import (
	"sync/atomic"

	"github.com/nasshu2916/dmx_viewer/pkg/logger"
)

// PanicHandler パニックハンドラー
type PanicHandler struct {
	logger *logger.Logger
	name   string
}

// NewPanicHandler パニックハンドラーを作成
func NewPanicHandler(logger *logger.Logger, name string) *PanicHandler {
	return &PanicHandler{
		logger: logger,
		name:   name,
	}
}

// Handle パニックを処理
func (ph *PanicHandler) Handle() {
	if r := recover(); r != nil {
		ph.logger.Error("Panic occurred", "component", ph.name, "panic", r)
	}
}

// CalculateUtilization チャンネル使用率計算
func CalculateUtilization(current, capacity int) float64 {
	if capacity == 0 {
		return 0.0
	}
	return float64(current) / float64(capacity) * 100.0
}

// DetermineHealthStatus 健全性ステータスを判定
func DetermineHealthStatus(utilization float64, droppedPackets int64) HealthStatus {
	if utilization >= CriticalUtilizationThreshold || droppedPackets > 0 {
		return CriticalStatus
	}
	if utilization >= HighUtilizationThreshold {
		return WarningStatus
	}
	return HealthyStatus
}

// LogChannelStats チャンネル統計をログ出力
func LogChannelStats(logger *logger.Logger, channelType ChannelType, queueLength, bufferSize int, droppedPackets int64, status HealthStatus) {
	utilization := CalculateUtilization(queueLength, bufferSize)

	fields := []interface{}{
		"channelType", channelType.String(),
		"queueLength", queueLength,
		"bufferSize", bufferSize,
		"utilization", utilization,
		"droppedPackets", droppedPackets,
		"status", status.String(),
	}

	switch status {
	case CriticalStatus:
		logger.Error("Critical channel utilization", fields...)
	case WarningStatus:
		logger.Warn("High channel utilization", fields...)
	default:
		logger.Debug("Channel statistics", fields...)
	}
}

// DropPacketWithLog パケットドロップをログ付きで実行
func DropPacketWithLog(logger *logger.Logger, counter *int64, channelType ChannelType, queueLength, bufferSize int, address string) {
	dropped := atomic.AddInt64(counter, 1)
	utilization := CalculateUtilization(queueLength, bufferSize)

	logger.Warn("Channel is full, dropping packet",
		"channelType", channelType.String(),
		"address", address,
		"droppedPackets", dropped,
		"queueLength", queueLength,
		"bufferSize", bufferSize,
		"utilization", utilization)
}
