package http

import (
	"fmt"
	"net/http"

	"github.com/nasshu2916/dmx_viewer/internal/infrastructure/artnet"
	"github.com/nasshu2916/dmx_viewer/internal/interface/httpctx"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
)

// MetricsHandler は Prometheus 互換のテキストフォーマットで簡易メトリクスを返す
// 外部依存を増やさずに /metrics を提供するための軽量実装。
type MetricsHandler struct {
	artnetServer *artnet.Server
	logger       *logger.Logger
}

func NewMetricsHandler(artnetServer *artnet.Server, logger *logger.Logger) *MetricsHandler {
	return &MetricsHandler{artnetServer: artnetServer, logger: logger}
}

// ServeHTTP implements http.Handler
func (h *MetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("metrics handler: ServeHTTP",
		"request_id", r.Header.Get("X-Request-Id"),
		"real_ip", httpctx.RealIP(r.Context()),
		"method", r.Method,
		"path", r.URL.Path,
	)

	bufferSize, recvQLen, sendQLen, droppedRecv, droppedSend := h.artnetServer.GetChannelStats()
	recvUtil, sendUtil := h.artnetServer.GetChannelUtilization()
	healthy, _ := h.artnetServer.IsChannelHealthy()

	healthValue := 0.0 // 0: healthy, 1: warning, 2: critical（大雑把に mapping）
	// 利用率 > 75% で warning、> 90% またはドロップ発生で critical
	if recvUtil > 90.0 || sendUtil > 90.0 || droppedRecv > 0 || droppedSend > 0 {
		healthValue = 2.0
	} else if recvUtil > 75.0 || sendUtil > 75.0 {
		healthValue = 1.0
	}

	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")

	// 各メトリクス（dmx_ プレフィクス）
	lines := []string{
		"# HELP dmx_artnet_channel_buffer_size ArtNet channel buffer capacity",
		"# TYPE dmx_artnet_channel_buffer_size gauge",
		fmt.Sprintf("dmx_artnet_channel_buffer_size %d", bufferSize),

		"# HELP dmx_artnet_receive_queue_length Number of items currently queued in receive channel",
		"# TYPE dmx_artnet_receive_queue_length gauge",
		fmt.Sprintf("dmx_artnet_receive_queue_length %d", recvQLen),

		"# HELP dmx_artnet_send_queue_length Number of items currently queued in send channel",
		"# TYPE dmx_artnet_send_queue_length gauge",
		fmt.Sprintf("dmx_artnet_send_queue_length %d", sendQLen),

		"# HELP dmx_artnet_dropped_receive_packets Dropped receive packets (current counter value)",
		"# TYPE dmx_artnet_dropped_receive_packets gauge",
		fmt.Sprintf("dmx_artnet_dropped_receive_packets %d", droppedRecv),

		"# HELP dmx_artnet_dropped_send_packets Dropped send packets (current counter value)",
		"# TYPE dmx_artnet_dropped_send_packets gauge",
		fmt.Sprintf("dmx_artnet_dropped_send_packets %d", droppedSend),

		"# HELP dmx_artnet_receive_utilization_percent Receive channel utilization percent",
		"# TYPE dmx_artnet_receive_utilization_percent gauge",
		fmt.Sprintf("dmx_artnet_receive_utilization_percent %.3f", recvUtil),

		"# HELP dmx_artnet_send_utilization_percent Send channel utilization percent",
		"# TYPE dmx_artnet_send_utilization_percent gauge",
		fmt.Sprintf("dmx_artnet_send_utilization_percent %.3f", sendUtil),

		"# HELP dmx_artnet_health_status 0=healthy,1=warning,2=critical (derived)",
		"# TYPE dmx_artnet_health_status gauge",
		fmt.Sprintf("dmx_artnet_health_status %.0f", healthValue),

		"# HELP dmx_artnet_overall_healthy 1 if healthy, else 0",
		"# TYPE dmx_artnet_overall_healthy gauge",
		fmt.Sprintf("dmx_artnet_overall_healthy %.0f", boolToFloat(healthy)),
	}

	for _, l := range lines {
		_, _ = w.Write([]byte(l + "\n"))
	}
}

func boolToFloat(b bool) float64 {
	if b {
		return 1.0
	}
	return 0.0
}
