package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/kcarretto/paragon/agent/transport"
	"go.uber.org/zap"

	"github.com/pkg/errors"
)

type request struct {
	Output []byte `json:"output"`
}

type task struct {
	ID      string `json:"id"`
	Content []byte `json"content"`
}

type response struct {
	Tasks []task `json:"tasks"`
}

type Transport struct {
	url    string
	logger *zap.Logger
	writer transport.TaskWriter
}

func (t Transport) Write(p []byte) (int, error) {
	resp, err := http.Post(t.url, "application/json", bytes.NewBuffer(p))
	if err != nil {
		return 0, errors.Wrap(err, "failed to connect via http")
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.logger.DPanic("Failed to read response data", zap.Error(err))
	}

	var content response
	if err := json.Unmarshal(data, &content); err != nil {
		t.logger.DPanic("response failed to unmarshal json", zap.Error(err))
	}

	for _, tk := range content.Tasks {
		t.writer.WriteTask(tk.ID, bytes.NewBuffer(tk.Content))
	}

	return len(p), nil
}

func (t Transport) Close() error {
	return nil
}
