package agent

import "time"

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
