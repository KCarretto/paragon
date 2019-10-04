package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/kcarretto/paragon/agent/transport"
	"go.uber.org/zap"

	"github.com/pkg/errors"
)

// New starts and returns a transport that receives tasks and reports output over HTTP.
func New(logger *zap.Logger, tasks transport.Tasker) (io.WriteCloser, error) {
	return Transport{
		url:    "http://127.0.0.1:8080",
		logger: logger,
		tasks:  tasks,
	}, nil
}

type Transport struct {
	url    string
	logger *zap.Logger
	tasks  transport.Tasker
}

func (t Transport) Write(payload []byte) (int, error) {
	// Send Agent Payload
	resp, err := http.Post(t.url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return 0, fmt.Errorf("failed to connect via http: %w", err)
	}
	defer resp.Body.Close()

	// Read Server Payload
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.logger.DPanic("Failed to read server payload!", zap.Error(err))
		return 0, fmt.Errorf("failed to read http response body: %w", err)
	}

	// Decode Server Payload
	var payload transport.ServerPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		t.logger.DPanic("Failed to unmarshal server payload!", zap.Error(err))
		return 0, errors.Wrap(err, "failed to unmarshal response from json")
	}

	t.logger.Debug("Received payload from server")

	// Queue Tasks
	for _, task := range payload.Tasks {
		t.tasks.QueueTask(task.ID, bytes.NewBuffer(task.Content))
	}

	return len(p), nil
}

func (t Transport) Close() error {
	return nil
}
