package router

import (
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: logging logic here
		next.ServeHTTP(w, r)
	})
}
