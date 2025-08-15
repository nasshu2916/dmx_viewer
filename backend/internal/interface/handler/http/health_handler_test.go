package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nasshu2916/dmx_viewer/internal/config"
	"github.com/nasshu2916/dmx_viewer/internal/infrastructure/artnet"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandler_Healthz_And_Readyz(t *testing.T) {
	l := logger.NewLogger("test")
	cfg := &config.ArtNet{PollIntervalSeconds: 300}
	server := artnet.NewServer(l, cfg)

	h := NewHealthHandler(server, l)

	// /healthz should be OK by default (no load, no drops)
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()
	h.Healthz(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// /readyz should be NOT READY because server is not running (no UDP listener)
	req2 := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rr2 := httptest.NewRecorder()
	h.Readyz(rr2, req2)
	assert.Equal(t, http.StatusServiceUnavailable, rr2.Code)
}
