package agent

import (
	"time"

	"github.com/kcarretto/paragon/agent/transport"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// An Option enables additional customization over an agent's configuration.
type Option func(*Agent)

// New initializes and configures a new Agent.
func New(logger *zap.Logger, taskExecutor Executor, options ...Option) Agent {
	agent := Agent{
		Tasks:         taskExecutor,
		Transports:    transport.Registry{},
		logger:        logger,
		logBufferSize: 1024 * 10,
	}

	for _, opt := range options {
		opt(&agent)
	}

	agent.queue = make(chan Task, agent.maxTaskBacklog)
	agent.buffer = transport.NewBuffer(make([]byte, 0, agent.logBufferSize))
	agent.logger = agent.logger.Named("agent").With(zap.Time("start_time", time.Now())).WithOptions(
		zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewTee(
				core,
				zapcore.NewCore(
					zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
					agent.buffer,
					zapcore.InfoLevel,
				))
		}),
	)

	agent.logger.Debug("Initialized agent")

	return agent
}

// SetTaskWorkers configures the number of worker routines used to run tasks.
func SetTaskWorkers(count int) Option {
	return func(agent *Agent) {
		agent.numWorkers = count
	}
}

// SetTaskBacklog configures the maximum number of tasks that can be backlogged before queue operations
// will start blocking.
func SetTaskBacklog(max int) Option {
	return func(agent *Agent) {
		agent.maxTaskBacklog = max
	}
}

// SetLogBufferSize configures the size (in bytes) of the initial log output buffer.
func SetLogBufferSize(size int) Option {
	return func(agent *Agent) {
		agent.logBufferSize = size
	}
}
