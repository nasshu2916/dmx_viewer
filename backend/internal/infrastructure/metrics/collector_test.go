package metrics

import (
	"testing"
	"time"

	"github.com/nasshu2916/dmx_viewer/internal/config"
	"github.com/nasshu2916/dmx_viewer/internal/infrastructure/artnet"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestCollector_ExportsReceivedPacketMetrics(t *testing.T) {
	log := logger.NewLogger("fatal")
	cfg := &config.ArtNet{PollIntervalSeconds: 300, ChannelBufferSize: 8}
	s := artnet.NewServer(log, cfg)

	// 記録: 現在5件、1秒前に6件、2秒前に7件、60秒前に2件、70秒前に3件
	now := time.Now().Truncate(time.Second)
	for i := 0; i < 5; i++ {
		s.RecordPacketAtForTest(now)
	}
	for i := 0; i < 6; i++ {
		s.RecordPacketAtForTest(now.Add(-1 * time.Second))
	}
	for i := 0; i < 7; i++ {
		s.RecordPacketAtForTest(now.Add(-2 * time.Second))
	}
	for i := 0; i < 2; i++ {
		s.RecordPacketAtForTest(now.Add(-60 * time.Second))
	}
	for i := 0; i < 3; i++ {
		s.RecordPacketAtForTest(now.Add(-70 * time.Second))
	}

	reg := prometheus.NewRegistry()
	_ = reg.Register(NewArtNetMetricsCollector(s))

	mfs, err := reg.Gather()
	assert.NoError(t, err)

	var total, lastMinute, lastSecond float64
	for _, mf := range mfs {
		if mf.GetName() == "dmx_artnet_received_packets_total" && len(mf.Metric) > 0 && mf.Metric[0].Counter != nil {
			total = mf.Metric[0].Counter.GetValue()
		}
		if mf.GetName() == "dmx_artnet_received_packets_last_minute" && len(mf.Metric) > 0 && mf.Metric[0].Gauge != nil {
			lastMinute = mf.Metric[0].Gauge.GetValue()
		}
		if mf.GetName() == "dmx_artnet_received_packets_last_second" && len(mf.Metric) > 0 && mf.Metric[0].Gauge != nil {
			lastSecond = mf.Metric[0].Gauge.GetValue()
		}
	}

	assert.Equal(t, float64(23), total)
	assert.Equal(t, float64(13), lastMinute)
	assert.Equal(t, float64(6), lastSecond)
}
