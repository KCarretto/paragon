package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kcarretto/paragon/transport"
	"go.uber.org/zap"
)

// Transport sends agent responses and receives server payloads via HTTP.
type Transport struct {
	transport.PayloadWriter
	URL    string
	Logger *zap.Logger
}

func (t Transport) Write(response []byte) (int, error) {
	// Send Agent Payload
	resp, err := http.Post(t.URL, "application/json", bytes.NewBuffer(response))
	if err != nil {
		return 0, fmt.Errorf("failed to connect via http: %w", err)
	}
	defer resp.Body.Close()

	// Read Server Payload
	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read http response body: %w", err)
	}

	if t.Logger == nil {
		t.Logger = zap.NewNop()
	}
	t.Logger.Debug("Received new payload from server", zap.Int("bytes", len(payload)))

	// Write payload
	t.WritePayload(payload)

	return len(response), nil
}

// Close is a no-op for an HTTP transport since it is not stateful.
func (t Transport) Close() error {
	return nil
}
