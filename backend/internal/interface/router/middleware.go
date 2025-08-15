package router

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/nasshu2916/dmx_viewer/internal/interface/httpctx"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
)

// コンテキストキー等は httpctx パッケージで管理

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

// RecovererMiddleware: panic を回収し 500 を返す
func RecovererMiddleware(l *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
			defer func() {
				if rec := recover(); rec != nil {
					l.Error("panic recovered",
						"panic", rec,
						"method", r.Method,
						"path", r.URL.Path,
						"remote_ip", httpctx.RealIP(r.Context()),
						"request_id", httpctx.RequestID(r.Context()),
						"stack", string(debug.Stack()),
					)
					http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(rw, r)
		})
	}
}

// RequestIDMiddleware: Request-ID を採番しヘッダ/コンテキストに格納
func RequestIDMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := httpctx.RequestIDFromHeaderOrNew(r)
			w.Header().Set("X-Request-Id", reqID)
			r = r.WithContext(httpctx.WithRequestID(r.Context(), reqID))
			next.ServeHTTP(w, r)
		})
	}
}

// RealIPMiddleware: Real-IP を解決しコンテキストに格納
func RealIPMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rip := httpctx.ResolveRealIP(r)
			r = r.WithContext(httpctx.WithRealIP(r.Context(), rip))
			next.ServeHTTP(w, r)
		})
	}
}

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
			next.ServeHTTP(w, r)
		})
	}
}
