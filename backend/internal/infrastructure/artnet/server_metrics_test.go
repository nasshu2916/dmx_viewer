package artnet

import (
	"testing"
	"time"

	"github.com/nasshu2916/dmx_viewer/internal/config"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func newTestServer() *Server {
	log := logger.NewLogger("fatal")
	cfg := &config.ArtNet{PollIntervalSeconds: 300, ChannelBufferSize: 8}
	return NewServer(log, cfg)
}

func TestReceivedPacketCounters_TotalAndLastMinute(t *testing.T) {
	s := newTestServer()

	now := time.Now().Truncate(time.Second)

	// 現在時刻のバケットに10件追加
	for i := 0; i < 10; i++ {
		s.recordPacketAt(now)
	}
	assert.Equal(t, int64(10), s.GetReceivedPacketsTotal())
	assert.Equal(t, int64(10), s.GetReceivedPacketsLastMinute())

	// 61秒前のバケットに5件追加（直近60秒からは除外される）
	past := now.Add(-61 * time.Second)
	for i := 0; i < 5; i++ {
		s.recordPacketAt(past)
	}
	// totalは加算、last minuteは変化なし
	assert.Equal(t, int64(15), s.GetReceivedPacketsTotal())
	assert.Equal(t, int64(10), s.GetReceivedPacketsLastMinute())
}

func TestReceivedPacketCounters_RollingWindowReset(t *testing.T) {
	s := newTestServer()

	// 120秒前の同一インデックスの古いデータを入れておく
	now := time.Now().Truncate(time.Second)
	veryOld := now.Add(-120 * time.Second)
	for i := 0; i < 7; i++ {
		s.recordPacketAt(veryOld)
	}

	// 現在時刻に1件追加すると、同じインデックスのバケットが現在秒に更新されリセットされる
	s.recordPacketAt(now)

	// 直近60秒の合計は現在の1件のみ（古い分は除外）
	assert.Equal(t, int64(8), s.GetReceivedPacketsTotal())
	assert.Equal(t, int64(1), s.GetReceivedPacketsLastMinute())
}
