package http

import (
	"io/fs"
	"net/http"

	"github.com/nasshu2916/dmx_viewer/internal/interface/httpctx"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
)

type StaticHandler struct {
	indexHtml []byte
	assetsFS  fs.FS
	logger    *logger.Logger
}

func NewStaticHandler(indexHtml []byte, assetsFS fs.FS, logger *logger.Logger) *StaticHandler {
	return &StaticHandler{
		indexHtml: indexHtml,
		assetsFS:  assetsFS,
		logger:    logger,
	}
}

func (h *StaticHandler) GetIndex(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("static handler: GetIndex",
		"request_id", r.Header.Get("X-Request-Id"),
		"real_ip", httpctx.RealIP(r.Context()),
		"method", r.Method,
		"path", r.URL.Path,
	)
	w.Header().Set("Content-Type", "text/html")
	w.Write(h.indexHtml)
}

func (h *StaticHandler) AssetsHandler() http.Handler {
	return http.StripPrefix("/assets/", http.FileServer(http.FS(h.assetsFS)))
}
