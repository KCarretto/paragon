package agent

import (
	"time"

	"go.uber.org/zap"
)

// An Option enables the configuration of an agent
type Option func(*agent)

// WithLogger configures the logger used by the agent.
func WithLogger(logger *zap.Logger) Option {
	return func(a *agent) {
		a.logger = logger
	}
}

// WithReporter configures the reporter used by the agent.
func WithReporter(reporter Reporter) Option {
	return func(a *agent) {
		a.reporter = reporter
	}
}

// WithInterval configures the interval the agent sleeps in between reports.
func WithInterval(interval time.Duration) Option {
	return func(a *agent) {
		a.interval = interval
	}
}

// WithJitter configures the bounds for a randomized duration that is added to the interval. The
// actual time an agent will sleep is in the range (interval - jitter, interval + jitter).
func WithJitter(jitter time.Duration) Option {
	return func(a *agent) {
		a.jitter = jitter
	}
}

// WithExecTimeout configures the default amount of time a task is permitted to run before signaling
// it to abort.
func WithExecTimeout(timeout time.Duration) Option {
	return func(a *agent) {
		a.execTimeout = timeout
	}
}
