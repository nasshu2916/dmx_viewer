package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nasshu2916/dmx_viewer/internal/interface/httpctx"
)

func TestRequestIDMiddleware_SetsHeaderAndContext(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := httpctx.RequestID(r.Context()); got == "" {
			t.Fatalf("request id not found in context")
		}
		w.WriteHeader(http.StatusOK)
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

func TestRequestIDMiddleware_PreservesExistingHeader(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Request-Id", "custom-id-1")

	RequestIDMiddleware()(next).ServeHTTP(rr, req)

	if got := rr.Header().Get("X-Request-Id"); got != "custom-id-1" {
		t.Fatalf("expected preserve header, got %q", got)
	}
}
