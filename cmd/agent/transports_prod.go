package main

import (
	"net/url"
	"os"

	"github.com/kcarretto/paragon/pkg/agent/transport"
	"github.com/kcarretto/paragon/pkg/agent/transport/http"
	"go.uber.org/zap"
)

func transports(logger *zap.Logger) (transports []transport.AgentMessageWriter) {
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

	transports = append(transports, &http.AgentTransport{
		URL: httpURL,
	})
	return
}
