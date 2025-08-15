package router

import (
	"net/http"
	"time"

	chimw "github.com/go-chi/chi/middleware"
)

// TimeoutMiddleware は chi の Timeout を薄くラップする
func TimeoutMiddleware(d time.Duration) func(http.Handler) http.Handler {
	return chimw.Timeout(d)
}
