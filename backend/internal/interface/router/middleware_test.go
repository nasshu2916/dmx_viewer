package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nasshu2916/dmx_viewer/internal/interface/httpctx"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
)

func TestRequestIDMiddleware_SetsHeaderAndContext(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := httpctx.RequestID(r.Context()); got == "" {
			t.Fatalf("request id not found in context")
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	RequestIDMiddleware()(next).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if rr.Header().Get("X-Request-Id") == "" {
		t.Fatalf("X-Request-Id header must be set")
	}
}

func TestRealIPMiddleware_XFFAndFallback(t *testing.T) {
	// X-Forwarded-For 優先
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := httpctx.RealIP(r.Context())
		if ip != "203.0.113.1" {
			t.Fatalf("expected real ip from XFF, got %s", ip)
		}
		w.WriteHeader(http.StatusOK)
	})
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Forwarded-For", "203.0.113.1, 70.41.3.18, 150.172.238.178")
	RealIPMiddleware()(next).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	// Fallback: RemoteAddr から
	next2 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := httpctx.RealIP(r.Context())
		if ip == "" {
			t.Fatalf("real ip should be set from RemoteAddr")
		}
		w.WriteHeader(http.StatusOK)
	})
	rr2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/", nil)
	RealIPMiddleware()(next2).ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr2.Code)
	}
}

func TestRecoverer_WithAccessLogOutside_Records500(t *testing.T) {
	l := logger.NewLogger("error")

	// panic を起こす最終ハンドラ
	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})

	// AccessLog(外) -> RequestID -> RealIP -> Recoverer(内) の順で合成
	h := AccessLogMiddleware(l)(RequestIDMiddleware()(RealIPMiddleware()(RecovererMiddleware(l)(final))))

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rr.Code)
	}
	if rr.Header().Get("X-Request-Id") == "" {
		t.Fatalf("X-Request-Id should be set even on panic")
	}
}
