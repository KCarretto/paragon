package event

import (
	"context"

	"github.com/kcarretto/paragon/ent"
)

// TaskQueuedEvent is a struct used to serialize a Task Queued event to pubsub
type TaskQueuedEvent struct {
	Target      *ent.Target
	Task        *ent.Task
	Credentials []*ent.Credential
	Tags        []*ent.Tag
}

// Publisher is a generic interface for publishing to a topic
type Publisher interface {
	Publish(context.Context, []byte) error
}
