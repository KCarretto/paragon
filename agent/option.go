package agent

import (
	"time"

	"github.com/kcarretto/paragon/agent/transport"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Option func(*Agent)

func New(logger *zap.Logger, runner Runner, options ...Option) Agent {
	agent := Agent{
		Publisher:  struct{}{}, // TODO: Default impl
		Subscriber: struct{}{}, // TODO: Default impl
		Scheduler:  struct{}{}, // TODO: Default impl
		Tasks:      runner,
		Transports: transport.Registry{},
		logger:     logger,
	}

	for _, opt := range options {
		opt(&agent)
	}

	agent.queue = make(chan Task, agent.maxTaskBacklog)
	agent.buffer = transport.NewBuffer(agent.maxLogBacklog)
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

func SetWorkers(count int) Option {
	return func(agent *Agent) {
		agent.numWorkers = count
	}
}

func SetTaskBacklog(max int) Option {
	return func(agent *Agent) {
		agent.maxTaskBacklog = max
	}
}

func SetLogBacklog(max int) Option {
	return func(agent *Agent) {
		agent.maxLogBacklog = max
	}
}
