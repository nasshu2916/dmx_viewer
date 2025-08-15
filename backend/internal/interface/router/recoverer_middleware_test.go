package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nasshu2916/dmx_viewer/pkg/logger"
)

func TestRecovererMiddleware_ConvertsPanicTo500(t *testing.T) {
	l := logger.NewLogger("error")
	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})
	h := RecovererMiddleware(l)(final)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rr.Code)
	}
	if rr.Body.String() == "" {
		t.Fatalf("expected body to be written")
	}
}
