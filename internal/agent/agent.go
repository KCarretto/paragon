package agent

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Agent interface {
	Execute(ctx context.Context, tasks ...Task)
	Run(ctx context.Context)

	ID() string
	String() string
	Interval() time.Duration
	Jitter() time.Duration
	ExecTimeout() time.Duration
}

type agent struct {
	logger   *zap.Logger
	reporter Reporter

	results chan Result
	pending sync.WaitGroup

	id               string
	interval         time.Duration
	jitter           time.Duration
	execTimeout      time.Duration
	pickTimeout      time.Duration
	maxResultBacklog uint16
}

func (a *agent) ID() string {
	return a.id
}

func (a *agent) String() string {
	return fmt.Sprintf("id: %s", a.ID())
}

func (a *agent) Interval() time.Duration {
	return a.interval
}

func (a *agent) Jitter() time.Duration {
	return a.jitter
}

func (a *agent) ExecTimeout() time.Duration {
	return a.execTimeout
}

func New(options ...Option) *agent {
	a := &agent{
		id:               uuid.New().String(),
		maxResultBacklog: 25,
		jitter:           time.Second * 10,
		interval:         time.Second * 60,
		execTimeout:      time.Second * 20,
		pickTimeout:      time.Second * 1,
		reporter:         &DebugReporter{},
		logger:           zap.NewNop(),
	}

	for _, opt := range options {
		opt(a)
	}

	a.results = make(chan Result, a.maxResultBacklog)

	a.logger.Debug(
		"Initialized new agent",
		zap.String("id", a.id),
		zap.Stringer("reporter", a.reporter),
		zap.Duration("interval", a.interval),
		zap.Duration("jitter", a.jitter),
		zap.Duration("exec_timeout", a.execTimeout),
	)

	return a
}
