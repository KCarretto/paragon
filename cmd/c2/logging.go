package main

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type respObserver struct {
	http.ResponseWriter
	Status int
	Body   []byte
}

func (resp *respObserver) Write(p []byte) (int, error) {
	resp.Body = append(resp.Body, p...)
	return resp.ResponseWriter.Write(p)
}

func (resp *respObserver) WriteHeader(code int) {
	resp.Status = code
	resp.ResponseWriter.WriteHeader(code)
}

func withLogging(logger *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		resp := respObserver{
			ResponseWriter: w,
			Status:         http.StatusOK}

		next.ServeHTTP(&resp, r)

		end := time.Now()

		log := logger.With(
			zap.Time("req_start", start),
			zap.Time("req_end", end),
			zap.Duration("req_latency", end.Sub(start)),
			zap.Int("req_status", resp.Status),
			zap.String("req_uri", r.RequestURI),
		)
		if resp.Status != http.StatusOK {
			log.Error("Failed to handle request", zap.Error(fmt.Errorf(string(resp.Body))))
		} else {
			log.Info("Successfully handled request")
		}
	})
}
