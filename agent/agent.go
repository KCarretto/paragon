package agent

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/kcarretto/paragon/agent/transport"
	"go.uber.org/zap"
)

// A Runner is responsible for executing tasks.
type Runner interface {
	Run(context.Context, *zap.Logger, Task)
}

// A Task contains instructions that will be executed by the runner, and a unique identifier.
type Task struct {
	ID      string
	Content io.Reader
}

// An Agent queues tasks for execution, executes them using the provided runner, and logs output to
// a registered transport.
type Agent struct {
	Tasks      Runner
	Transports transport.Registry

	numWorkers     int
	maxTaskBacklog int
	maxLogBacklog  int

	wg     sync.WaitGroup
	buffer *transport.Buffer
	logger *zap.Logger
	queue  chan Task
}

func (agent Agent) assertReady() {
	if agent.Tasks == nil ||
		agent.buffer == nil ||
		agent.logger == nil ||
		agent.queue == nil {

		panic(ErrAgentUninitialized)
	}
}

// Run the agent, enabling tasks to be queued and output to be logged to a registered transport.
func (agent Agent) Run() {
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

	// Send buffer to a registered transport
	for {
		if err := agent.send(agent.logger.Named("writer"), agent.buffer); err != nil {
			agent.logger.DPanic("Failed to send buffer", zap.Error(err))
		}
	}
}

// QueueTask implements transport.Tasker by initializing a task and adding it to the execution queue.
func (agent Agent) QueueTask(id string, content io.Reader) {
	agent.assertReady()

	task := Task{
		id,
		content,
	}
	agent.queue <- task
}

// Close the agent by finishing all tasks in the queue.
func (agent Agent) Close() error {
	if agent.queue != nil {
		close(agent.queue)
	}
	agent.wg.Wait()
	return nil
}

// taskWorker consumes tasks from the queue and executes them using the configured runner.
func (agent Agent) taskWorker(ctx context.Context, logger *zap.Logger, queue <-chan Task) {
	for task := range queue {
		tLogger := logger.With(zap.String("task_id", task.ID))
		tLogger.Debug("Starting task execution")
		defer tLogger.Debug("Finished task execution")
		agent.Tasks.Run(ctx, tLogger, task)
	}
}

// send the provided buffer using a registered transport. The transport is selected based on the
// ordering of the List() method in the transport registry. Send iterates through all registered
// transports until one succeeds. If no configured transport is successful, send returns an error.
func (agent Agent) send(logger *zap.Logger, buffer *transport.Buffer) error {
	for _, meta := range agent.Transports.List() {
		tLogger := logger.Named("transport").With(zap.String("transport_name", meta.Name))
		writer, err := agent.Transports.Get(meta.Name, tLogger, agent)
		if err != nil {
			tLogger.Error("Failed to get transport from registry", zap.Error(err))
			continue
		}
		if writer == nil {
			tLogger.Error("Failed to get transport from registry", zap.Error(err))
			continue
		}

		delay := meta.Interval - time.Since(buffer.Timestamp())
		time.Sleep(delay)

		if _, err = buffer.WriteTo(writer); err != nil {
			tLogger.Error("Failed to write buffer to transport", zap.Error(err))
			// TODO: Remove, just decrement priority, or just close?
			if err := agent.Transports.Remove(meta.Name); err != nil {
				tLogger.Error("Failed to close transport", zap.Error(err))
			}
		}
	}

	return ErrNoTransportAvailable
}
