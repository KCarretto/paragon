package agent

import (
	"context"
)

// Tasks can be run, which return a result
type Task interface {
	ID() string
	Run(context.Context, Config) Result
}
