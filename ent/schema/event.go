package schema

import (
	"time"

	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
)

// Event holds the schema definition for the Event entity.
type Event struct {
	ent.Schema
}

// EventKind is not great for readability but we did this so we do not have to define the Enum
// values twice. Below in the Kind field we need to pack the values in, this provides a way to do
// that by making an array.
type EventKind string

var eventKinds []string

func newEventKind(value string) EventKind {
	eventKinds = append(eventKinds, value)
	return EventKind(value)
}

// IF you update this you MUST regent the ent stuff
var (
	EventCreateJob EventKind = newEventKind("CREATE_JOB")
	EventCreateTag EventKind = newEventKind("CREATE_TAG")

	EventApplyTagToTask   EventKind = newEventKind("APPLY_TAG_TO_TASK")
	EventApplyTagToTarget EventKind = newEventKind("APPLY_TAG_TO_TARGET")
	EventApplyTagToJob    EventKind = newEventKind("APPLY_TAG_TO_JOB")

	EventRemoveTagFromTask   EventKind = newEventKind("REMOVE_TAG_FROM_TASK")
	EventRemoveTagFromTarget EventKind = newEventKind("REMOVE_TAG_FROM_TARGET")
	EventRemoveTagFromJob    EventKind = newEventKind("REMOVE_TAG_FROM_JOB")

	EventCreateTarget           EventKind = newEventKind("CREATE_TARGET")
	EventSetTargetFields        EventKind = newEventKind("SET_TARGET_FIELDS")
	EventDeleteTarget           EventKind = newEventKind("DELETE_TARGET")
	EventAddCredentialForTarget EventKind = newEventKind("ADD_CREDENTIAL_FOR_TARGET")

	EventUploadFile EventKind = newEventKind("UPLOAD_FILE")

	EventCreateLink    EventKind = newEventKind("CREATE_LINK")
	EventSetLinkFields EventKind = newEventKind("SET_LINK_FIELDS")

	EventActivateUser EventKind = newEventKind("ACTIVATE_USER")
	EventCreateUser   EventKind = newEventKind("CREATE_USER")
	EventMakeAdmin    EventKind = newEventKind("MAKE_ADMIN")
	EventRemoveAdmin  EventKind = newEventKind("REMOVE_ADMIN")
	EventChangeName   EventKind = newEventKind("CHANGE_NAME")

	EventActivateService EventKind = newEventKind("ACTIVATE_SERVICE")
	EventCreateService   EventKind = newEventKind("CREATE_SERVICE")

	EventLikeEvent EventKind = newEventKind("LIKE_EVENT")

	EventOther EventKind = newEventKind("OTHER")
)

// Fields of the Event.
func (Event) Fields() []ent.Field {
	return []ent.Field{
		field.Time("CreationTime").
			Default(func() time.Time {
				return time.Now()
			}).
			Comment("The timestamp for when the Job was created"),
		field.Enum("Kind").
			Values(eventKinds...).
			Comment("The kind of event"),
	}
}

// Edges of the Event.
func (Event) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("job", Job.Type).
			Unique().
			Comment("An Event can have a job"),
		edge.To("file", File.Type).
			Unique().
			Comment("An Event can have a file"),
		edge.To("credential", Credential.Type).
			Unique().
			Comment("An Event can have a credential"),
		edge.To("link", Link.Type).
			Unique().
			Comment("An Event can have a link"),
		edge.To("tag", Tag.Type).
			Unique().
			Comment("An Event can have a tag"),
		edge.To("target", Target.Type).
			Unique().
			Comment("An Event can have a target"),
		edge.To("task", Task.Type).
			Unique().
			Comment("An Event can have a task"),
		edge.To("user", User.Type).
			Unique().
			Comment("An Event can have a user"),
		edge.To("event", Event.Type).
			Unique().
			Comment("An Event can have an event"),
		edge.To("service", Service.Type).
			Unique().
			Comment("An Event can have a service"),
		edge.To("likers", User.Type).
			Comment("An Event can have a few likers"),
		edge.From("owner", User.Type).
			Ref("events").
			Unique().
			Required().
			Comment("An Job must have an owner"),
	}
}
