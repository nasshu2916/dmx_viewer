package router

import (
	"net/http"

	"github.com/nasshu2916/dmx_viewer/internal/interface/httpctx"
)

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
