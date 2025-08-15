package http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nasshu2916/dmx_viewer/internal/config"
	"github.com/nasshu2916/dmx_viewer/internal/infrastructure/artnet"
	metrics "github.com/nasshu2916/dmx_viewer/internal/infrastructure/metrics"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestMetricsHandler_ServeHTTP(t *testing.T) {
	l := logger.NewLogger("test")
	cfg := &config.ArtNet{PollIntervalSeconds: 300}
	server := artnet.NewServer(l, cfg)

	// カスタムRegistryを構築
	reg := metrics.BuildRegistry(server)
	mh := NewMetricsHandlerWithRegistry(reg, l)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rr := httptest.NewRecorder()
	mh.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	body, _ := io.ReadAll(rr.Body)
	s := string(body)
	assert.Contains(t, s, "dmx_artnet_channel_buffer_size")
	assert.Contains(t, s, "dmx_artnet_overall_healthy")
}
