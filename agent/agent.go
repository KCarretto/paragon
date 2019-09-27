package agent

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/kcarretto/paragon/agent/transport"
	"go.uber.org/zap"
)

// A Executor is responsible for executing tasks.
type Executor interface {
	Exec(context.Context, *zap.Logger, Task)
}

// A Task contains instructions that will be executed by the runner, and a unique identifier.
type Task struct {
	ID      string
	Content io.Reader
}

// An Agent queues tasks for execution, executes them, and logs output to a registered transport.
type Agent struct {
	Tasks      Executor
	Transports *transport.Registry

	numWorkers     int
	maxTaskBacklog int
	logBufferSize  int

	wg     sync.WaitGroup
	buffer *transport.Buffer
	logger *zap.Logger
	queue  chan Task
}

// assertReady ensures that the agent is properly initialized, otherwise panics.
func (agent Agent) assertReady() {
	if agent.Tasks == nil ||
		agent.buffer == nil ||
		agent.logger == nil ||
		agent.queue == nil {

		panic(ErrAgentUninitialized)
	}
}

// Run the agent, enabling tasks to be queued and output to be logged to a registered transport.
func (agent *Agent) Run() error {
	agent.assertReady()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start workers for running tasks.
	for i := 0; i < agent.numWorkers || i < 1; i++ {
		agent.wg.Add(1)
		go func() {
			defer func() {
				agent.logger.Debug("Stopped task worker")
				agent.wg.Done()
			}()
			agent.logger.Debug("Started task worker")
			agent.taskWorker(ctx, agent.logger.Named("tasks").With(zap.Int("worker_id", i)), agent.queue)
		}()
	}

	// Send buffer to a registered transport.
	for {
		if err := agent.send(agent.logger.Named("transport"), agent.buffer); err != nil {
			agent.logger.Error("Failed to send buffer", zap.Error(err))
			return err
		}
	}
}

// QueueTask implements transport.Tasker by initializing a task and adding it to the execution queue.
func (agent *Agent) QueueTask(id string, content io.Reader) {
	agent.assertReady()

	task := Task{
		id,
		content,
	}
	agent.queue <- task
}

// Close the agent by finishing all tasks in the queue.
func (agent *Agent) Close() error {
	if agent.queue != nil {
		close(agent.queue)
	}
	agent.wg.Wait()
	// TODO: Close transports in registry
	return nil
}

// taskWorker consumes tasks from the queue and executes them using the configured runner.
func (agent *Agent) taskWorker(ctx context.Context, logger *zap.Logger, queue <-chan Task) {
	for task := range queue {
		tLogger := logger.With(zap.String("task_id", task.ID))
		tLogger.Debug("Starting task execution")
		defer tLogger.Debug("Finished task execution")
		agent.Tasks.Exec(ctx, tLogger, task)
	}
}

// send the provided buffer using a registered transport. The transport is selected based on the
// ordering of the List() method in the transport registry. Send iterates through all registered
// transports until one succeeds. If no configured transport is successful, send returns an error.
func (agent *Agent) send(logger *zap.Logger, buffer *transport.Buffer) error {

	for _, meta := range agent.Transports.List() {
		tLogger := logger.Named(meta.Name)
		writer, err := agent.Transports.Get(meta.Name, tLogger, agent)
		if err != nil {
			tLogger.Error("Failed to get transport from registry", zap.Error(err))
			continue
		}

		delay := meta.Interval - time.Since(buffer.Timestamp())
		time.Sleep(delay)

		n, err := buffer.WriteTo(writer)
		if err != nil && err != io.EOF {
			// TODO: Decrement priority, or just close?
			closeErr := agent.Transports.CloseTransport(meta.Name)
			tLogger.Error(
				"Failed to write buffer to transport",
				zap.Error(err),
				zap.NamedError("transport_close_err", closeErr),
			)
			continue
		}

		if n > 0 {
			tLogger.Debug("Successfully transported output", zap.Int64("bytes_sent", n))
		}

		return nil
	}

	return ErrNoTransportAvailable
}
