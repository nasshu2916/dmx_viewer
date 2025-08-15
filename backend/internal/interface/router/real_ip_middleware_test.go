package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nasshu2916/dmx_viewer/internal/interface/httpctx"
)

func TestRealIPMiddleware_XForwardedForPreferred(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := httpctx.RealIP(r.Context()); got != "203.0.113.1" {
			t.Fatalf("expected XFF ip, got %q", got)
		}
		w.WriteHeader(http.StatusOK)
	})
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Forwarded-For", "203.0.113.1, 70.41.3.18, 150.172.238.178")
	RealIPMiddleware()(next).ServeHTTP(rr, req)
}

func TestRealIPMiddleware_XRealIPFallback(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := httpctx.RealIP(r.Context()); got != "198.51.100.5" {
			t.Fatalf("expected X-Real-IP, got %q", got)
		}
		w.WriteHeader(http.StatusOK)
	})
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Real-IP", "198.51.100.5")
	RealIPMiddleware()(next).ServeHTTP(rr, req)
}
