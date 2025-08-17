package http

import (
	"net/http"

	"github.com/nasshu2916/dmx_viewer/internal/interface/httpctx"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsHandler は promhttp に委譲してメトリクスを返す（カスタムRegistry対応）
type MetricsHandler struct {
	handler http.Handler
	logger  *logger.Logger
}

// NewMetricsHandlerWithRegistry は指定された Registry を利用してハンドラを返す
func NewMetricsHandlerWithRegistry(reg prometheus.Gatherer, logger *logger.Logger) *MetricsHandler {
	return &MetricsHandler{handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}), logger: logger}
}

// ServeHTTP implements http.Handler
func (h *MetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("metrics handler: ServeHTTP",
		"request_id", r.Header.Get("X-Request-Id"),
		"real_ip", httpctx.RealIP(r.Context()),
		"method", r.Method,
		"path", r.URL.Path,
	)
	h.handler.ServeHTTP(w, r)
}
