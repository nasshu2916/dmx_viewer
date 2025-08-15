package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	httpHandler "github.com/nasshu2916/dmx_viewer/internal/interface/handler/http"
	"github.com/nasshu2916/dmx_viewer/internal/interface/handler/websocket"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
)

func NewRouter(static *httpHandler.StaticHandler, timeHandler *httpHandler.TimeHandler, health *httpHandler.HealthHandler, metrics *httpHandler.MetricsHandler, ws *websocket.WebSocketHandler, l *logger.Logger, httpTimeout time.Duration) http.Handler {
	r := chi.NewRouter()

	// ベース（全体）ミドルウェア
	// 順序: AccessLog(外側) -> RequestID -> RealIP -> (各グループ固有) -> Recoverer(内側)
	r.Use(AccessLogMiddleware(l))
	r.Use(RequestIDMiddleware())
	r.Use(RealIPMiddleware())

	// 静的/HTTP API グループ（タイムアウト適用）
	r.Group(func(gr chi.Router) {
		gr.Use(ForceTimeoutMiddleware(httpTimeout))
		gr.Use(RecovererMiddleware(l))

		gr.Get("/", static.GetIndex)
		gr.Handle("/assets/*", static.AssetsHandler())
		gr.Get("/api/time", timeHandler.GetTime)
		gr.Get("/healthz", health.Healthz)
		gr.Get("/readyz", health.Readyz)
		gr.Handle("/metrics", metrics)
	})

	// WebSocket グループ（タイムアウトは適用しない）
	r.Group(func(gr chi.Router) {
		gr.Use(RecovererMiddleware(l))
		gr.Handle("/ws", http.HandlerFunc(ws.ServeWS))
	})

	return r
}
