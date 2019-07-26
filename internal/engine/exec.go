package engine

import "context"

// Result stores and displays the result of running a task.
type Result struct {
	TaskID    string
	TimeStart int
	TimeEnd   int
	Stdout    string
	Stderr    string
	ExitCode  int
	Err       bool
}

// Tasks can be run, which return a result
type Task interface {
	Run(ctx context.Context) Result
}

func (a *agent) Execute(ctx context.Context, tasks ...Task) {
	for _, t := range tasks {
		// TODO: Log task scheduling
		a.pending.Add(1)
		go func(task Task) {
			// TODO: Enable script timeout
			r := task.Run(ctx)
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
