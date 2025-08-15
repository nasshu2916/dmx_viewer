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

	r.Use(LoggingMiddleware(l))
	r.Get("/", static.GetIndex)
	r.Handle("/assets/*", static.AssetsHandler())

	r.Get("/api/time", time.GetTime)

	r.Handle("/ws", http.HandlerFunc(ws.ServeWS))

	return r
}
