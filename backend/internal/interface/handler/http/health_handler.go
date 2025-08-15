package http

import (
	"encoding/json"
	"net/http"

	"github.com/nasshu2916/dmx_viewer/internal/infrastructure/artnet"
	"github.com/nasshu2916/dmx_viewer/internal/interface/httpctx"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
)

type HealthHandler struct {
	artnetServer *artnet.Server
	logger       *logger.Logger
}

func NewHealthHandler(artnetServer *artnet.Server, logger *logger.Logger) *HealthHandler {
	return &HealthHandler{
		artnetServer: artnetServer,
		logger:       logger,
	}
}

// /healthz — チャネル健全性の確認（Queue使用率・ドロップ）
func (h *HealthHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("health handler: Healthz",
		"request_id", r.Header.Get("X-Request-Id"),
		"real_ip", httpctx.RealIP(r.Context()),
		"method", r.Method,
		"path", r.URL.Path,
	)

	healthy, msg := h.artnetServer.IsChannelHealthy()
	status := http.StatusOK
	resp := map[string]interface{}{
		"status":  "ok",
		"message": "",
	}
	if !healthy {
		status = http.StatusServiceUnavailable
		resp["status"] = "unhealthy"
		resp["message"] = msg
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}

// /readyz — ArtNetサーバーのリスナ確立確認
func (h *HealthHandler) Readyz(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("health handler: Readyz",
		"request_id", r.Header.Get("X-Request-Id"),
		"real_ip", httpctx.RealIP(r.Context()),
		"method", r.Method,
		"path", r.URL.Path,
	)

	ready := h.artnetServer.IsRunning()
	status := http.StatusOK
	resp := map[string]interface{}{
		"status": "ready",
	}
	if !ready {
		status = http.StatusServiceUnavailable
		resp["status"] = "not_ready"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}
