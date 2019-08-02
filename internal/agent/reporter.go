package agent

import (
	"context"
	"fmt"
)

// A Reporter is responsible for reporting results to a C2, and returning new tasks that need to be
// scheduled for execution.
type Reporter interface {
	String() string
	Report(context.Context, Agent, []Result) []Task
}

// DebugReporter is a default reporter that simply prints results to stdout and does not
// fetch new tasks.
type DebugReporter struct{}

func (reporter *DebugReporter) String() string {
	return "debug_reporter"
}

// Report enables the DebugReporter to print results to stdout.
func (reporter *DebugReporter) Report(ctx context.Context, agent Agent, results []Result) []Task {
	for _, result := range results {
		fmt.Printf("Result: %+v\n", result)
	}

	return nil
}
