package engine

import (
	"fmt"
)

// Reporter TODO
type Reporter interface {
	Report([]Result) []Task
}

type DebugReporter struct{}

func (reporter *DebugReporter) Report(results []Result) []Task {
	for _, result := range results {
		fmt.Printf("Result: %+v\n", result)
	}

	return nil
}
