package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type observer struct {
	http.ResponseWriter
	StatusCode int
}

// WriteHeader wraps the underlying ResponseWriter to preserve the response status code.
func (o *observer) WriteHeader(statusCode int) {
	o.StatusCode = statusCode
	o.ResponseWriter.WriteHeader(statusCode)
}

// WithLogging provides request logging for an http handler.
func WithLogging(log *zap.Logger) func(http.Handler) http.HandlerFunc {
	return func(f http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrapped := &observer{w, http.StatusOK}
			f.ServeHTTP(wrapped, r)

			msg := log.With(
				zap.Duration("duration", time.Since(start)),
				zap.String("method", r.Method),
				zap.String("uri", r.RequestURI),
				zap.Int64("length", r.ContentLength),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()),
				zap.Int("status", wrapped.StatusCode),
			)

			if wrapped.StatusCode != http.StatusOK {
				msg.Error("Failed to handle request")
			} else {
				msg.Info("Successfully handled request")
			}
		}
	}
}
