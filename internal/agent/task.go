package agent

import (
	"context"
)

// A Task can be run, which return a result that will be reported.
type Task interface {
	ID() string
	Run(context.Context, Config) Result
}
