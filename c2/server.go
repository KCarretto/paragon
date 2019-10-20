package c2

import (
	"github.com/kcarretto/paragon/api/events"
)

// Server manages communication with agents.
type Server struct {
	*Queue

	OnTaskExecuted func(events.TaskExecuted)
}
