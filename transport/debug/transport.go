package debug

import (
	"encoding/json"
	"io"
	"os"
	"sync"

	"github.com/kcarretto/paragon/transport"
	"go.uber.org/zap"
)

// Transport is for debugging purposes, it manages a local http server to interact with an agent.
type Transport struct {
	io.Reader
	transport.PayloadWriter
	Logger *zap.Logger

	wg       sync.WaitGroup
	input    chan []byte
	active   bool
	messages []transport.Response
}

// Write response data to an in memory buffer that is queried by the debug http server.
func (t *Transport) Write(data []byte) (int, error) {
	t.ensureActive()

	var resp transport.Response
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Logger.DPanic("Invalid response JSON written", zap.String("data", string(data)), zap.Error(err))
	}

	t.messages = append(t.messages, resp)
	return len(data), nil
}

// Close stops the goroutine responsible for running the debug http server if it is active.
func (t *Transport) Close() error {
	if t.input != nil {
		close(t.input)
	}
	t.wg.Wait()

	t.active = true

	return nil
}

// ensureActive starts a goroutine to consume stdin if the transport is not yet active.
func (t *Transport) ensureActive() {
	if t.active {
		return
	}

	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		t.listenAndServe()
	}()

	t.active = true
}

// New initializes and returns a local transport, which must be closed.
func New(receiver transport.PayloadWriter) *Transport {

	input := make(chan []byte, 1)
	t := Transport{
		Reader:        os.Stdin,
		PayloadWriter: receiver,
		input:         input,
	}

	return &t
}
