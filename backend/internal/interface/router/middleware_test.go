package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nasshu2916/dmx_viewer/internal/interface/httpctx"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
)

func TestLoggingMiddleware_SetsRequestID(t *testing.T) {
	l := logger.NewLogger("error")

	// 次ハンドラ: 200で応答
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ミドルウェアで設定された Request-ID/Real-IP が参照できること
		if got := httpctx.RequestID(r.Context()); got == "" {
			t.Fatalf("request id not found in context")
		}
		// Real-IP は RemoteAddr 由来のIPが入る（テスト環境のデフォルトを許容）
		_ = httpctx.RealIP(r.Context())
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	mw := LoggingMiddleware(l)
	mw(next).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	got := rr.Header().Get("X-Request-Id")
	if got == "" {
		t.Fatalf("X-Request-Id header must be set")
	}

	// ヘッダに指定があればそれを尊重
	rr = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Request-Id", "test-id-123")
	mw(next).ServeHTTP(rr, req)
	if rr.Header().Get("X-Request-Id") != "test-id-123" {
		t.Fatalf("X-Request-Id should be preserved from request header")
	}
}

func TestLoggingMiddleware_PanicRecovery(t *testing.T) {
	l := logger.NewLogger("error")

	// 次ハンドラ: panic を起こす
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/panic", nil)

	mw := LoggingMiddleware(l)
	// recoverされ、500が返ること
	mw(next).ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rr.Code)
	}

	if rr.Header().Get("X-Request-Id") == "" {
		t.Fatalf("X-Request-Id should be set even on panic")
	}
}
