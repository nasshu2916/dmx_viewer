package router

import (
	"net/http"
	"time"
)

// ForceTimeoutMiddleware は http.TimeoutHandler を利用して強制的にタイムアウト応答を返す
func ForceTimeoutMiddleware(d time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.TimeoutHandler(next, d, "request timeout")
	}
}
