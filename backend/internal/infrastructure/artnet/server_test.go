package artnet

import (
	"testing"

	"github.com/nasshu2916/dmx_viewer/internal/config"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestServer_ChannelBuffering(t *testing.T) {
	logger := logger.NewLogger("test")
	cfg := &config.ArtNet{
		PollIntervalSeconds: 30, // テスト中はポーリングを無効にする
		ChannelBufferSize:   10, // 小さなバッファサイズでテスト
	}

	server := NewServer(logger, cfg)

	// 初期状態の確認
	bufferSize, queueLength, droppedPackets := server.GetChannelStats()
	assert.Equal(t, 10, bufferSize)
	assert.Equal(t, 0, queueLength)
	assert.Equal(t, int64(0), droppedPackets)
}

func TestServer_DefaultConfiguration(t *testing.T) {
	logger := logger.NewLogger("test")
	cfg := &config.ArtNet{
		PollIntervalSeconds: 30,
	}

	server := NewServer(logger, cfg)

	// デフォルト値の確認
	bufferSize, queueLength, droppedPackets := server.GetChannelStats()
	assert.Equal(t, DefaultChannelBufferSize, bufferSize)
	assert.Equal(t, 0, queueLength)
	assert.Equal(t, int64(0), droppedPackets)
}

func TestServer_StatsReset(t *testing.T) {
	logger := logger.NewLogger("test")
	cfg := &config.ArtNet{
		PollIntervalSeconds: 30,
		ChannelBufferSize:   5,
	}

	server := NewServer(logger, cfg)

	// 初期状態
	assert.Equal(t, int64(0), server.GetDroppedPackets())

	// ドロップパケット数をリセット
	server.ResetDroppedPackets()
	assert.Equal(t, int64(0), server.GetDroppedPackets())
}

func TestServer_ChannelCapacity(t *testing.T) {
	logger := logger.NewLogger("test")
	cfg := &config.ArtNet{
		PollIntervalSeconds: 30,
		ChannelBufferSize:   2,
	}

	// 非常に小さなバッファでテスト
	server := NewServer(logger, cfg)

	// チャネルが適切なサイズで作成されているか確認
	bufferSize, _, _ := server.GetChannelStats()
	assert.Equal(t, 2, bufferSize)

	// チャネルの容量確認
	packetChan := server.PacketChan()
	assert.Equal(t, 2, cap(packetChan))
}
