package http

import (
	"net/http"

	"github.com/nasshu2916/dmx_viewer/internal/interface/httpctx"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsHandler は promhttp に委譲してメトリクスを返す
type MetricsHandler struct {
	handler http.Handler
	logger  *logger.Logger
}

func NewMetricsHandler(_ interface{}, logger *logger.Logger) *MetricsHandler {
	return &MetricsHandler{handler: promhttp.Handler(), logger: logger}
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
