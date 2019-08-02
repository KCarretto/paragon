package agent

import "time"

// Config describes the getter methods for the agent's configuration, and additionally provides
// the Reconfigure() method to apply configuration updates at runtime.
type Config interface {
	Reconfigure(updates ...Option)

	ID() string
	String() string
	Interval() time.Duration
	Jitter() time.Duration
	ExecTimeout() time.Duration
}

func (a *agent) Reconfigure(updates ...Option) {
	for _, opt := range updates {
		opt(a)
	}
}
