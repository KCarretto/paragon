package event

import (
	"github.com/kcarretto/paragon/ent"
)


// TaskQueued is a struct used to serialize a Task Queued event to pubsub
type TaskQueued struct {
	Target      *ent.Target
	Task        *ent.Task
	Credentials []*ent.Credential
	Tags        []*ent.Tag
}
