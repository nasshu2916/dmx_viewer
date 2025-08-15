package metrics

import (
	"errors"
	"github.com/nasshu2916/dmx_viewer/internal/infrastructure/artnet"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

// ArtNetMetricsCollector は ArtNet サーバの状態を収集する Prometheus Collector
type ArtNetMetricsCollector struct {
	server *artnet.Server

	bufferSizeDesc     *prometheus.Desc
	recvQLenDesc       *prometheus.Desc
	sendQLenDesc       *prometheus.Desc
	droppedRecvDesc    *prometheus.Desc
	droppedSendDesc    *prometheus.Desc
	recvUtilPercent    *prometheus.Desc
	sendUtilPercent    *prometheus.Desc
	healthStatusDesc   *prometheus.Desc
	overallHealthyDesc *prometheus.Desc
}

func NewArtNetMetricsCollector(server *artnet.Server) *ArtNetMetricsCollector {
	return &ArtNetMetricsCollector{
		server: server,
		bufferSizeDesc: prometheus.NewDesc(
			"dmx_artnet_channel_buffer_size",
			"ArtNet channel buffer capacity",
			nil, nil,
		),
		recvQLenDesc: prometheus.NewDesc(
			"dmx_artnet_receive_queue_length",
			"Number of items currently queued in receive channel",
			nil, nil,
		),
		sendQLenDesc: prometheus.NewDesc(
			"dmx_artnet_send_queue_length",
			"Number of items currently queued in send channel",
			nil, nil,
		),
		droppedRecvDesc: prometheus.NewDesc(
			"dmx_artnet_dropped_receive_packets",
			"Dropped receive packets (current counter value)",
			nil, nil,
		),
		droppedSendDesc: prometheus.NewDesc(
			"dmx_artnet_dropped_send_packets",
			"Dropped send packets (current counter value)",
			nil, nil,
		),
		recvUtilPercent: prometheus.NewDesc(
			"dmx_artnet_receive_utilization_percent",
			"Receive channel utilization percent",
			nil, nil,
		),
		sendUtilPercent: prometheus.NewDesc(
			"dmx_artnet_send_utilization_percent",
			"Send channel utilization percent",
			nil, nil,
		),
		healthStatusDesc: prometheus.NewDesc(
			"dmx_artnet_health_status",
			"0=healthy,1=warning,2=critical (derived)",
			nil, nil,
		),
		overallHealthyDesc: prometheus.NewDesc(
			"dmx_artnet_overall_healthy",
			"1 if healthy, else 0",
			nil, nil,
		),
	}
}

func (c *ArtNetMetricsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.bufferSizeDesc
	ch <- c.recvQLenDesc
	ch <- c.sendQLenDesc
	ch <- c.droppedRecvDesc
	ch <- c.droppedSendDesc
	ch <- c.recvUtilPercent
	ch <- c.sendUtilPercent
	ch <- c.healthStatusDesc
	ch <- c.overallHealthyDesc
}

func (c *ArtNetMetricsCollector) Collect(ch chan<- prometheus.Metric) {
	bufferSize, recvQLen, sendQLen, droppedRecv, droppedSend := c.server.GetChannelStats()
	recvUtil, sendUtil := c.server.GetChannelUtilization()
	healthy, _ := c.server.IsChannelHealthy()

	healthValue := 0.0
	if recvUtil > 90.0 || sendUtil > 90.0 || droppedRecv > 0 || droppedSend > 0 {
		healthValue = 2.0
	} else if recvUtil > 75.0 || sendUtil > 75.0 {
		healthValue = 1.0
	}

	ch <- prometheus.MustNewConstMetric(c.bufferSizeDesc, prometheus.GaugeValue, float64(bufferSize))
	ch <- prometheus.MustNewConstMetric(c.recvQLenDesc, prometheus.GaugeValue, float64(recvQLen))
	ch <- prometheus.MustNewConstMetric(c.sendQLenDesc, prometheus.GaugeValue, float64(sendQLen))
	ch <- prometheus.MustNewConstMetric(c.droppedRecvDesc, prometheus.GaugeValue, float64(droppedRecv))
	ch <- prometheus.MustNewConstMetric(c.droppedSendDesc, prometheus.GaugeValue, float64(droppedSend))
	ch <- prometheus.MustNewConstMetric(c.recvUtilPercent, prometheus.GaugeValue, recvUtil)
	ch <- prometheus.MustNewConstMetric(c.sendUtilPercent, prometheus.GaugeValue, sendUtil)
	ch <- prometheus.MustNewConstMetric(c.healthStatusDesc, prometheus.GaugeValue, healthValue)
	if healthy {
		ch <- prometheus.MustNewConstMetric(c.overallHealthyDesc, prometheus.GaugeValue, 1)
	} else {
		ch <- prometheus.MustNewConstMetric(c.overallHealthyDesc, prometheus.GaugeValue, 0)
	}
}

// RegisterArtNetMetrics はカスタム Collector と標準 Collector を登録する
func RegisterArtNetMetrics(server *artnet.Server, _ *logger.Logger) {
	// デフォルトレジストリにプロセス/Go/ビルド情報 Collector を登録（重複登録は無視）
	_ = registerIfNeeded(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	_ = registerIfNeeded(collectors.NewGoCollector())
	_ = registerIfNeeded(collectors.NewBuildInfoCollector())

	// ArtNetカスタムメトリクス（重複登録は無視）
	_ = registerIfNeeded(NewArtNetMetricsCollector(server))
}

func registerIfNeeded(c prometheus.Collector) error {
	if err := prometheus.Register(c); err != nil {
		var are prometheus.AlreadyRegisteredError
		if errors.As(err, &are) {
			return nil
		}
		return err
	}
	return nil
}
