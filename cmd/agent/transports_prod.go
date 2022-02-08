package main

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"

	"github.com/kcarretto/paragon/pkg/agent/transport"
	httplib "github.com/kcarretto/paragon/pkg/agent/transport/http"
	"go.uber.org/zap"
)

func transports(logger *zap.Logger) (transports []transport.AgentMessageWriter) {
	transports = append(transports,
		&httplib.AgentTransport{
			URL: &url.URL{
				Scheme: "http",
				Host:   "localhost:4444",
			},
			Client: &http.Client{
				Timeout: time.Second * 5,
				Transport: &http.Transport{
					TLSHandshakeTimeout:   time.Second * 5,
					DisableKeepAlives:     true,
					ResponseHeaderTimeout: time.Second * 5,
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
					},
				},
			},
		},
	)
	return
}
