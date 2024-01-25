package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// LogMiddleware log the request
func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			next.ServeHTTP(w, r)
			log.Info().Msgf("%s [%s] %s in %v", r.RemoteAddr, r.Method, r.URL.Path, time.Since(start))
		}()
	})
}
