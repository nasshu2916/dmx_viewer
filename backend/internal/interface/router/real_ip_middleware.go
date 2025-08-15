package router

import (
	"net/http"

	"github.com/nasshu2916/dmx_viewer/internal/interface/httpctx"
)

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
