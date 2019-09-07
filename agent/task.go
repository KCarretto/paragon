package agent

import (
	"context"
	"io"
	"time"

	"github.com/kcarretto/paragon/agent/transport"
	"go.uber.org/zap"
)

type Task struct {
	Publisher
	Subscriber
	Scheduler
	Transports transport.Registrar

	ID      string
	Content io.Reader
}

func (agent Agent) QueueTask(id string, content io.Reader) {
	agent.assertReady()

	task := Task{
		agent.Publisher,
		agent.Subscriber,
		agent.Scheduler,
		agent.Transports,

		id,
		content,
	}
	agent.queue <- task
}

func (agent Agent) worker(ctx context.Context, logger *zap.Logger, queue <-chan Task) {
	for task := range queue {
		taskLogger := logger.With(zap.String("task_id", task.ID))
		taskLogger.Debug("Running task")
		agent.Tasks.Run(ctx, taskLogger, task)
	}
}

func (agent Agent) send(logger *zap.Logger, buffer *transport.Buffer) error {
	for _, meta := range agent.Transports.List() {
		tLogger := logger.With(zap.String("transport_name", meta.Name))

		writer, err := agent.Transports.Get(meta.Name, tLogger.Named("transport"), agent)
		if err != nil {
			// TODO: Log error
			continue
		}
		if writer == nil {
			// TODO: Log error
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
