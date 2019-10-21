package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/kcarretto/paragon/agent"
)

// Sender is an HTTP Sender.
type Sender struct {
	URL *url.URL
	*http.Client
}

// Send an agent message to the server using http.
func (s Sender) Send(w agent.ServerMessageWriter, msg agent.Message) error {
	// Encode agent message
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to encode msg: %w", err)
	}

	// Select destination url
	url := "http://127.0.0.1:8080/"
	if s.URL != nil {
		url = s.URL.String()
	}

	// Create http request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Select http client
	var client *http.Client = s.Client
	if client == nil {
		client = http.DefaultClient
	}

	// Sent http request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send msg: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid http response code: %d", resp.StatusCode)
	}

	// Read http response
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read http response body: %w", err)
	}

	// Decode http response
	var srvMsg agent.ServerMessage
	if err := json.Unmarshal(content, &srvMsg); err != nil {
		return fmt.Errorf("failed to unmarshal http response body: %w", err)
	}

	// Write server message
	w.WriteServerMessage(&srvMsg)

	return nil
}
