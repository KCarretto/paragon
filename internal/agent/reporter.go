package agent

import (
	"context"
	"fmt"
)

// Reporter TODO
type Reporter interface {
	String() string
	Report(context.Context, Agent, []Result) []Task
}

type DebugReporter struct{}

func (reporter *DebugReporter) String() string {
	return "debug_reporter"
}

func (reporter *DebugReporter) Report(ctx context.Context, agent Agent, results []Result) []Task {
	for _, result := range results {
		fmt.Printf("Result: %+v\n", result)
	}

	return nil
}
