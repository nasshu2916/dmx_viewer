package router

import (
	"net/http"
	"time"

	"github.com/nasshu2916/dmx_viewer/internal/interface/httpctx"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
)

// AccessLogMiddleware: メソッド/パス/ステータス/所要時間/Real-IP/Request-ID を出力
func AccessLogMiddleware(l *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
			defer func() {
				duration := time.Since(start)
				// Recoverer より外側で動作するため、Request-ID はレスポンスヘッダから取得
				reqID := rw.Header().Get("X-Request-Id")
				// Real-IP はヘッダから都度解決
				realIP := httpctx.ResolveRealIP(r)
				l.Info("http request",
					"method", r.Method,
					"path", r.URL.Path,
					"status", rw.status,
					"bytes", rw.bytes,
					"duration_ms", duration.Milliseconds(),
					"remote_ip", realIP,
					"request_id", reqID,
				)
			}()
			next.ServeHTTP(rw, r)
		})
	}
}
