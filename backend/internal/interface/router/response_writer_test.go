package router

import (
	"bufio"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResponseWriter_CapturesStatusAndBytes(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := &responseWriter{ResponseWriter: rec, status: http.StatusOK}

	rw.WriteHeader(http.StatusCreated)
	if got, want := rw.status, http.StatusCreated; got != want {
		t.Fatalf("status not captured: got %d want %d", got, want)
	}
	if rec.Code != http.StatusCreated {
		t.Fatalf("response code not forwarded: got %d", rec.Code)
	}

	n, err := rw.Write([]byte("abc"))
	if err != nil || n != 3 {
		t.Fatalf("write failed or bytes mismatch: n=%d err=%v", n, err)
	}
	if rw.bytes != 3 {
		t.Fatalf("bytes not captured: got %d", rw.bytes)
	}
	if rec.Body.String() != "abc" {
		t.Fatalf("body not forwarded: %q", rec.Body.String())
	}
}

// --- fakes ---
type fakeRW struct {
	http.ResponseWriter
	hijacked bool
	flushed  bool
	pushed   []string
}

func (f *fakeRW) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	f.hijacked = true
	return nil, nil, nil
}

func (f *fakeRW) Flush() { f.flushed = true }

func (f *fakeRW) Push(target string, _ *http.PushOptions) error {
	f.pushed = append(f.pushed, target)
	return nil
}

func TestResponseWriter_DelegatesFlush(t *testing.T) {
	base := httptest.NewRecorder()
	f := &fakeRW{ResponseWriter: base}
	rw := &responseWriter{ResponseWriter: f, status: http.StatusOK}
	rw.Flush()
	if !f.flushed {
		t.Fatalf("expected underlying Flush to be called")
	}
}

func TestResponseWriter_DelegatesHijack_Supported(t *testing.T) {
	base := httptest.NewRecorder()
	f := &fakeRW{ResponseWriter: base}
	rw := &responseWriter{ResponseWriter: f, status: http.StatusOK}
	_, _, err := rw.Hijack()
	if err != nil {
		t.Fatalf("unexpected error from Hijack: %v", err)
	}
	if !f.hijacked {
		t.Fatalf("expected underlying Hijack to be called")
	}
}

func TestResponseWriter_DelegatesHijack_NotSupported(t *testing.T) {
	base := httptest.NewRecorder() // does not implement Hijacker
	rw := &responseWriter{ResponseWriter: base, status: http.StatusOK}
	if _, _, err := rw.Hijack(); err == nil {
		t.Fatalf("expected error when underlying does not support Hijacker")
	}
}

func TestResponseWriter_DelegatesPush(t *testing.T) {
	base := httptest.NewRecorder()
	f := &fakeRW{ResponseWriter: base}
	rw := &responseWriter{ResponseWriter: f, status: http.StatusOK}
	if err := rw.Push("/asset.js", nil); err != nil {
		t.Fatalf("unexpected error from Push: %v", err)
	}
	if len(f.pushed) != 1 || f.pushed[0] != "/asset.js" {
		t.Fatalf("expected push to be delegated, got %v", f.pushed)
	}
}

func TestResponseWriter_DelegatesPush_NotSupported(t *testing.T) {
	base := httptest.NewRecorder() // does not implement Pusher
	rw := &responseWriter{ResponseWriter: base, status: http.StatusOK}
	if err := rw.Push("/asset.js", nil); err == nil {
		t.Fatalf("expected error when underlying does not support Pusher")
	}
}
