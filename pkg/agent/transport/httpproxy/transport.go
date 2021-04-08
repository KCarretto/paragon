package httpproxy

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/kcarretto/paragon/pkg/agent/transport"
	"go.uber.org/zap"
)

// AgentTransport provides an HTTP transport for the agent to communicate with a server.
type AgentTransport struct {
	transport.Transport

	URL *url.URL
	*http.Client

	PROXY string
}

// WriteAgentMessage sends an agent message to the server using HTTP(s).
func (t *AgentTransport) WriteAgentMessage(ctx context.Context, w transport.ServerMessageWriter, msg transport.AgentMessage) error {
	// Encode agent message
	msgBuffer := new(bytes.Buffer)
	if err := t.EncodeAgentMessage(msg, msgBuffer); err != nil {
		return err
	}

	//creating the proxyURL
	proxyStr := "http://localhost:3128"
	if t.PROXY != "" {
		proxyStr = t.PROXY
	}

	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		return fmt.Errorf("failed to create http proxy: %w", err)
	}

	//adding the proxy settings to the Transport object
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	// Select destination url
	url := "http://127.0.0.1:8080/"
	if t.URL != nil {
		url = t.URL.String()
	}

	// Create http request
	req, err := http.NewRequest(http.MethodPost, url, msgBuffer)
	if err != nil {
		return fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	// Select http client
	//adding the Transport object to the http Client
	var client *http.Client = t.Client
	if client == nil {
		client = &http.Client{Transport: transport}
	}

	// Send http request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send msg: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid http response code: %d", resp.StatusCode)
	}

	// Decode http response
	srvMsg, err := t.DecodeServerMessage(resp.Body)
	if err != nil {
		return err
	}

	// Write server message
	w.WriteServerMessage(ctx, srvMsg)

	return nil
}

// ServerTransport provides a transport for the server to handle HTTP communications from an agent.
type ServerTransport struct {
	transport.Transport
	Log    *zap.Logger
	Server transport.AgentMessageWriter
}

// ListenAndServe will begin serving the http transport on the provided address.
func (t *ServerTransport) ListenAndServe(addr string) error {
	t.Log.Info("Started HTTP transport", zap.String("listen_on", addr))
	return http.ListenAndServe(addr, t)
}

// ServeHTTP enables the transport to handle HTTP requests provided to the server.
func (t *ServerTransport) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Setup request logging
	logger := t.Log.With(
		zap.String("RemoteIP", req.RemoteAddr),
		zap.String("HostHeader", req.Header.Get("Host")),
		zap.Int64("ContentLength", req.ContentLength),
	)

	// Decode agent message
	msg, err := t.DecodeAgentMessage(req.Body)
	if err != nil {
		logger.Error("Failed to decode agent request", zap.Error(err))
		http.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	// Log Request
	logger = logger.With(
		zap.String("AgentID", msg.Metadata.AgentID),
		zap.String("SessionID", msg.Metadata.SessionID),
		zap.String("PrimaryIP", msg.Metadata.PrimaryIP),
		zap.String("PrimaryMAC", msg.Metadata.PrimaryMAC),
		zap.String("Hostname", msg.Metadata.Hostname),
		zap.String("UUID", msg.Metadata.MachineUUID),
	)
	logger.Info("Handling agent message", zap.Int("ResultCount", len(msg.Results)))

	// Prepare response
	resp := &response{
		t.Transport,
		t.Log,
		w,
	}

	// Provide the agent message to the server for handling
	if err := t.Server.WriteAgentMessage(req.Context(), resp, msg); err != nil {
		logger.Error("Failed to handle agent request", zap.Error(err))
	}
}

// response is a helper implementation of transport.ServerMessageWriter to enable writing server
// messages to an http.ResponseWriter.
type response struct {
	transport.Transport
	Log *zap.Logger

	w http.ResponseWriter
}

func (resp *response) WriteServerMessage(ctx context.Context, msg transport.ServerMessage) {
	if err := resp.EncodeServerMessage(msg, resp.w); err != nil {
		resp.Log.Error("failed to encode response to agent", zap.Error(err))
		http.Error(resp.w, "failed to encode response", http.StatusInternalServerError)
		return
	}
	resp.w.Header().Set("Content-Type", "application/json")
}
