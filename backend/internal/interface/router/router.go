package router

import (
	"net/http"

	"github.com/go-chi/chi"
	httpHandler "github.com/nasshu2916/dmx_viewer/internal/interface/handler/http"
	"github.com/nasshu2916/dmx_viewer/internal/interface/handler/websocket"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
)

func NewRouter(static *httpHandler.StaticHandler, time *httpHandler.TimeHandler, ws *websocket.WebSocketHandler, l *logger.Logger) http.Handler {
	r := chi.NewRouter()

	// 単機能ミドルウェアの組み合わせ
	// 順序: AccessLog(外側) -> RequestID -> RealIP -> Recoverer(内側)
	// これにより panic 時も Recoverer が 500 を書いた後に AccessLog が正しい status を記録できる
	r.Use(AccessLogMiddleware(l))
	r.Use(RequestIDMiddleware())
	r.Use(RealIPMiddleware())
	r.Use(RecovererMiddleware(l))
	r.Get("/", static.GetIndex)
	r.Handle("/assets/*", static.AssetsHandler())

	r.Get("/api/time", time.GetTime)

	r.Handle("/ws", http.HandlerFunc(ws.ServeWS))

	return r
}
