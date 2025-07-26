package artnet

import (
	"testing"
	"time"

	"github.com/nasshu2916/dmx_viewer/internal/config"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestChannelPressure(t *testing.T) {
	logger := logger.NewLogger("test")
	cfg := &config.ArtNet{
		PollIntervalSeconds: 300,
		ChannelBufferSize:   10,
	}

	server := NewServer(logger, cfg)

	// 初期状態の確認
	healthy, msg := server.IsChannelHealthy()
	assert.True(t, healthy)
	assert.Equal(t, "Channels are healthy", msg)

	// チャンネル使用率の確認
	receiveUtil, sendUtil := server.GetChannelUtilization()
	assert.Equal(t, 0.0, receiveUtil)
	assert.Equal(t, 0.0, sendUtil)
}

func TestChannelPressureSimulation(t *testing.T) {
	logger := logger.NewLogger("test")
	cfg := &config.ArtNet{
		PollIntervalSeconds: 300,
		ChannelBufferSize:   5,
	}

	server := NewServer(logger, cfg)
	pressureTest := NewChannelPressureTest(server)

	stats := pressureTest.CheckChannelPressure()
	assert.True(t, stats["isHealthy"].(bool))

	// 送信チャンネルを満杯にする
	dummyData := make([]byte, 100)
	for i := 0; i < 6; i++ { // バッファサイズより多く送信
		err := server.SendToWriteChanForTest(dummyData, &DummyAddr{})
		if i < 5 {
			assert.NoError(t, err, "Should succeed for first 5 packets")
		} else {
			assert.Error(t, err, "Should fail when buffer is full")
		}
	}

	// チャンネルが不健全になったことを確認
	healthy, _ := server.IsChannelHealthy()
	assert.False(t, healthy, "Channel should be unhealthy after dropping packets")

	// 使用率の確認
	_, sendUtil := server.GetChannelUtilization()
	assert.True(t, sendUtil > 75, "Send utilization should be high")
}

func TestChannelUtilizationLevels(t *testing.T) {
	logger := logger.NewLogger("test")
	cfg := &config.ArtNet{
		PollIntervalSeconds: 300,
		ChannelBufferSize:   10,
	}

	server := NewServer(logger, cfg)

	// 75%使用率のテスト
	dummyData := make([]byte, 100)
	for i := 0; i < 8; i++ { // 80%使用率
		err := server.SendToWriteChanForTest(dummyData, &DummyAddr{})
		assert.NoError(t, err)
	}

	healthy, msg := server.IsChannelHealthy()
	assert.False(t, healthy)
	assert.Contains(t, msg, "High utilization")

	// 統計情報のリセット
	server.ResetDroppedPackets()

	// チャンネルを空にする（実際のシナリオでは、runSenderが処理する）
	for len(server.sendChan) > 0 {
		<-server.sendChan
	}

	healthy, msg = server.IsChannelHealthy()
	assert.True(t, healthy)
	assert.Equal(t, "Channels are healthy", msg)
}

func TestWriteChanWithTimeout(t *testing.T) {
	logger := logger.NewLogger("test")
	cfg := &config.ArtNet{
		PollIntervalSeconds: 300,
		ChannelBufferSize:   2,
	}

	server := NewServer(logger, cfg)
	dummyData := make([]byte, 100)

	// チャンネルを満杯にする
	for i := 0; i < 2; i++ {
		err := server.SendToWriteChanForTest(dummyData, &DummyAddr{})
		assert.NoError(t, err)
	}

	// タイムアウト付き送信のテスト
	start := time.Now()
	err := server.SendToWriteChanWithTimeoutForTest(dummyData, &DummyAddr{}, 100*time.Millisecond)
	elapsed := time.Since(start)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "timeout")
	assert.True(t, elapsed >= 100*time.Millisecond)
	assert.True(t, elapsed < 200*time.Millisecond) // 余裕をもって
}
