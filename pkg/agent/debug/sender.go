package debug

import (
	"context"
	"net/http"
	"sync"

	"github.com/kcarretto/paragon/pkg/agent"
	"github.com/kcarretto/paragon/proto/codec"
	"go.uber.org/zap"
)

// Sender is for debugging purposes, it manages a local http server to interact with an agent.
type Sender struct {
	Log *zap.Logger

	srv      *http.Server
	active   bool
	wg       sync.WaitGroup
	messages []agent.Message
	tasks    []*codec.Task
}

// Send appends an agent message to the result array, and writes any queued tasks to the provided
// writer.
func (transport *Sender) Send(w agent.ServerMessageWriter, msg agent.Message) error {
	transport.ensureActive()

	transport.messages = append(transport.messages, msg)

	tasks := make([]*codec.Task, len(transport.tasks))
	copy(tasks, transport.tasks)
	transport.tasks = transport.tasks[:0]

	w.WriteServerMessage(&agent.ServerMessage{
		Tasks: tasks,
	})

	return nil
}

// QueueTask buffers a task for execution.
func (transport *Sender) QueueTask(task *codec.Task) {
	transport.tasks = append(transport.tasks, task)
}

// Close stops the goroutine responsible for running the debug http server if it is active.
func (transport *Sender) Close() (err error) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if transport.srv != nil {
		err = transport.srv.Shutdown(ctx)
	}

	transport.wg.Wait()
	transport.active = false

	return err
}

// ensureActive starts a goroutine to consume stdin if the transport is not yet active.
func (transport *Sender) ensureActive() {
	if transport.active {
		return
	}

	transport.wg.Add(1)
	go func() {
		defer transport.wg.Done()
		transport.listenAndServe()
	}()

	transport.active = true
}
