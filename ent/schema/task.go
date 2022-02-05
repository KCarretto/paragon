package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Task holds the schema definition for the Task entity.
type Task struct {
	ent.Schema
}

// Fields of the Task.
func (Task) Fields() []ent.Field {
	return []ent.Field{
		field.Time("QueueTime").
			Default(func() time.Time {
				return time.Now()
			}).
			Comment("The timestamp for when the Task was queued/created"),
		field.Time("LastChangedTime").
			Comment("The timestamp for when the Task was last changed"),
		field.Time("ClaimTime").
			Optional().
			Comment("The timestamp for when the Task was claim"),
		field.Time("ExecStartTime").
			Optional().
			Comment("The timestamp for when the Task was executed"),
		field.Time("ExecStopTime").
			Optional().
			Comment("The timestamp for when the Task's execution ended"),
		field.Text("Content").
			NotEmpty().
			Comment("The content of the task (usually a Renegade Script)"),
		field.Text("Output").
			Optional().
			Comment("The output from executing the task"),
		field.String("Error").
			Optional().
			Comment("The error, if any, produced while executing the Task"),
		field.String("SessionID").
			MaxLen(250).
			Optional(),
	}
}

// Edges of the Task.
func (Task) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tags", Tag.Type).
			Comment("A Task can have many Tags"),
		edge.From("job", Job.Type).
			Ref("tasks").
			Unique().
			Required().
			Comment("A Task must have a job"),
		edge.From("target", Target.Type).
			Ref("tasks").
			Unique().
			Comment("A Task might have a target"),
	}
}
