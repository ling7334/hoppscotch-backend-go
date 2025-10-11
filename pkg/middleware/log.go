package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// LogMiddleware log the request
func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			next.ServeHTTP(w, r)
			slog.Info(fmt.Sprintf("%s [%s] %s in %v", r.RemoteAddr, r.Method, r.URL.Path, time.Since(start)))
		}()
	})
}
