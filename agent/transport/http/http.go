package http

import (
	"bytes"
	"encoding/json"
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

func (t Transport) Write(p []byte) (int, error) {
	// Encode Agent Payload
	output, err := json.Marshal(
		transport.AgentPayload{
			Output: p,
		},
	)
	if err != nil {
		return 0, errors.Wrap(err, "failed to marshal output to json")
	}

	// Send Agent Payload
	resp, err := http.Post(t.url, "application/json", bytes.NewBuffer(output))
	if err != nil {
		return 0, errors.Wrap(err, "failed to connect via http")
	}
	defer resp.Body.Close()

	// Read Server Payload
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.logger.DPanic("Failed to read server payload!", zap.Error(err))
		return 0, errors.Wrap(err, "failed to read response body")
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
