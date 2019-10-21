// +build dev

package main

import (
	"net/url"
	"os"
	"time"

	"github.com/kcarretto/paragon/agent"
	"github.com/kcarretto/paragon/agent/http"
	"go.uber.org/zap"
)

func transports(logger *zap.Logger) (transports []agent.Transport) {
	httpURL := &url.URL{
		Scheme: "http",
		Host:   "127.0.0.1:8080",
	}
	if addr := os.Getenv("C2_HTTP_ADDR"); addr != "" {
		url, err := url.Parse(addr)
		if err != nil || url == nil {
			logger.Error("Failed to parse URL %q", zap.String("url", addr), zap.Error(err))
		} else {
			httpURL = url
		}
	}

	transports = append(transports, agent.Transport{
		Sender:   http.Sender{URL: httpURL},
		Log:      logger.Named("http").With(zap.String("http_url", httpURL.String())),
		Name:     "http",
		Interval: time.Second * 10,
		Jitter:   time.Second * 5,
	},
	)

	return
}
