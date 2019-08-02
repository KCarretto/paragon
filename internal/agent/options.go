package agent

import (
	"time"

	"go.uber.org/zap"
)

type Option func(*agent)

func WithLogger(logger *zap.Logger) Option {
	return func(a *agent) {
		a.logger = logger
	}
}

func WithReporter(reporter Reporter) Option {
	return func(a *agent) {
		a.reporter = reporter
	}
}

func WithInterval(interval time.Duration) Option {
	return func(a *agent) {
		a.interval = interval
	}
}

func WithJitter(jitter time.Duration) Option {
	return func(a *agent) {
		a.jitter = jitter
	}
}

func WithExecTimeout(timeout time.Duration) Option {
	return func(a *agent) {
		a.execTimeout = timeout
	}
}
