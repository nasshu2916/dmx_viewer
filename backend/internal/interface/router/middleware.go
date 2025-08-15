package router

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/nasshu2916/dmx_viewer/internal/interface/httpctx"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
)

// コンテキストキー等は httpctx パッケージで管理

// LoggingMiddleware は以下を行うミドルウェア:
// - リクエストIDの付与とレスポンスヘッダ出力
// - Real-IP の解決
// - panic リカバリと 500 応答
// - メソッド/パス/ステータス/所要時間のログ出力
func LoggingMiddleware(l *logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			reqID := httpctx.RequestIDFromHeaderOrNew(r)
			realIP := httpctx.ResolveRealIP(r)

			// レスポンスヘッダにリクエストIDを設定
			w.Header().Set("X-Request-Id", reqID)

			// ステータスとサイズを取得するためのラッパ
			rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

			// コンテキストに付与
			ctx := httpctx.WithRequestID(r.Context(), reqID)
			ctx = httpctx.WithRealIP(ctx, realIP)
			r = r.WithContext(ctx)

			defer func() {
				// panic リカバリ
				if rec := recover(); rec != nil {
					l.Error("panic recovered",
						"panic", rec,
						"method", r.Method,
						"path", r.URL.Path,
						"remote_ip", realIP,
						"request_id", reqID,
						"stack", string(debug.Stack()),
					)
					http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}

				// リクエスト完了ログ
				duration := time.Since(start)
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

			next.ServeHTTP(w, r)
		})
	}
}

// responseWriter はステータスコードと送信バイト数を記録するためのラッパ
type responseWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (w *responseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.bytes += n
	return n, err
}
