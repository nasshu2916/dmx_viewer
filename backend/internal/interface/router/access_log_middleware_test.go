package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nasshu2916/dmx_viewer/pkg/logger"
)

func TestAccessLogMiddleware_PreservesStatusAndBody(t *testing.T) {
	l := logger.NewLogger("error")

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte("ok"))
	})
	h := AccessLogMiddleware(l)(next)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", rr.Code)
	}
	if rr.Body.String() != "ok" {
		t.Fatalf("expected body 'ok', got %q", rr.Body.String())
	}
}
