package agent

import (
	"context"
	"io"
	"sync"

	"github.com/kcarretto/paragon/agent/transport"
	"go.uber.org/zap"
)

// Runner is responsible for executing tasks.
type Runner interface {
	Run(ctx context.Context, logger *zap.Logger, content io.Reader) error
}

type Agent struct {
	Runner

	wg         sync.WaitGroup
	output     chan []byte
	transports transport.Registry
}

func (agent Agent) Start(ctx context.Context) {
	// TODO: Configure output chan size
	// TODO: Should use buffer + mutex instead of chan []byte?
	agent.output = make(chan []byte, 100)
	agent.transports = transport.Registry{}

	go agent.writer(ctx)
}

// WriteTask concurrently runs a task.
func (agent Agent) WriteTask(id string, content io.Reader) {
	// TODO: Schedule instead of just running go function?
	agent.wg.Add(1)
	go func() {
		defer agent.wg.Done()
		// TODO: Context cancellation & configure logger
		if err := agent.Run(context.Background(), zap.NewNop(), content); err != nil {
			// TODO: Log error
		}
	}()
}

// Write output using an agent transport. Safe for concurrent use.
func (agent Agent) Write(p []byte) (int, error) {
	agent.output <- p
	return len(p), nil
}

// Close the agent by closing all open transports and finishing all executing tasks.
func (agent Agent) Close() error {
	defer agent.wg.Wait()
	return agent.transports.Close()
}
