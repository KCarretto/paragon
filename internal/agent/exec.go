package agent

import (
	"context"

	"go.uber.org/zap"
)

// An Executor is capable of executing tasks and scheduling their results to be reported.
type Executor interface {
	Execute(ctx context.Context, tasks ...Task)
}

// Execute schedules the execution of the provided tasks.
func (a *agent) Execute(ctx context.Context, tasks ...Task) {
	for _, t := range tasks {
		a.logger.Debug("Scheduling task execution", zap.String("task_id", t.ID()))
		a.pending.Add(1)
		go func(task Task) {
			// TODO: Enable task specific timeout
			r := task.Run(ctx, a)
			a.results <- r
			a.pending.Done()
		}(t)
	}
}

// pickResults returns a list of all Results that are currently available.
func (a *agent) pickResults(ctx context.Context) []Result {
	var results []Result
	for {
		select {
		case <-ctx.Done():
			return results
		case r := <-a.results:
			results = append(results, r)
		default:
			return results
		}
	}
}
