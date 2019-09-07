package agent

import (
	"context"
	"errors"
	"sync"

	"github.com/kcarretto/paragon/agent/transport"
	"go.uber.org/zap"
)

type Runner interface {
	Run(context.Context, *zap.Logger, Task)
}
type Publisher interface{}
type Subscriber interface{}
type Scheduler interface{}

type Agent struct {
	Publisher
	Subscriber
	Scheduler
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
		agent.Publisher == nil ||
		agent.Subscriber == nil ||
		agent.Scheduler == nil ||
		agent.buffer == nil ||
		agent.logger == nil ||
		agent.queue == nil {

		panic(errors.New("must initialize agent with agent.New()"))
	}
}

func (agent Agent) Run() {
	agent.assertReady()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < agent.numWorkers || i < 1; i++ {
		agent.wg.Add(1)
		go func() {
			defer func() {
				agent.logger.Debug("Stopped task worker")
				agent.wg.Done()
			}()
			agent.logger.Debug("Started task worker")
			agent.worker(ctx, agent.logger.Named("tasks").With(zap.Int("worker_id", i)), agent.queue)
		}()
	}

	for {
		if err := agent.send(agent.logger.Named("writer"), agent.buffer); err != nil {
			agent.logger.DPanic("Failed to send buffer", zap.Error(err))
		}
	}
}

func (agent Agent) Close() error {
	if agent.queue != nil {
		close(agent.queue)
	}
	agent.wg.Wait()
	return nil
}
