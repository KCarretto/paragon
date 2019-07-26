package engine

import (
	"context"
	"sync"
	"time"
)

// TODO: Agent interface & New method
type Option func(*agent)

type agent struct {
	results chan Result
	pending sync.WaitGroup

	interval         time.Duration
	minDelta         time.Duration
	maxDelta         time.Duration
	execTimeout      time.Duration
	pickTimeout      time.Duration
	maxResultBacklog uint16
	reporter         Reporter
}

func (a *agent) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			// TODO: Log
			a.pending.Wait()
			return
		default:
			// Collect completed task results
			pickCtx, cancel := context.WithTimeout(ctx, a.pickTimeout)
			defer cancel()
			results := a.pickResults(pickCtx)

			// Report results & fetch new tasks
			if a.reporter == nil {
				// TODO: Error handling
				continue
			}
			tasks := a.reporter.Report(results)

			// Schedule execution of new tasks
			execCtx, cancel := context.WithTimeout(ctx, a.execTimeout)
			defer cancel()
			a.Execute(execCtx, tasks...)

			// TODO: Log cycle complete

			// Sleep for interval, still permitting tasks to execute
			// TODO: Allow ctx to be cancelled even when sleeping
			a.delay()
		}
	}
}

func (a *agent) delay() {
	// TODO: Randomize interval with delta
	time.Sleep(a.interval)
}

func New(options ...Option) *agent {
	a := &agent{
		maxResultBacklog: 25,
		interval:         time.Second * 60,
		execTimeout:      time.Second * 20,
		pickTimeout:      time.Second * 1,
		reporter:         &DebugReporter{},
	}

	for _, opt := range options {
		opt(a)
	}

	// TODO: Configure result buffer size
	a.results = make(chan Result, a.maxResultBacklog)

	return a
}
