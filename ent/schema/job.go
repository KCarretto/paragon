package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Job holds the schema definition for the Job entity.
type Job struct {
	ent.Schema
}

// Fields of the Target.
func (Job) Fields() []ent.Field {
	return []ent.Field{
		field.String("Name").
			NotEmpty().
			Comment("The name of the job"),
		field.Time("CreationTime").
			Default(func() time.Time {
				return time.Now()
			}).
			Comment("The timestamp for when the Job was created"),
		field.Text("Content").
			NotEmpty().
			Comment("The content of the job (usually a Renegade Script)"),
		field.Bool("Staged").
			Comment("The boolean that represents if a job's tasks shall be emitted/returned from claimTasks (false means yes)"),
	}
}

// Edges of the Job.
func (Job) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tasks", Task.Type).
			Comment("A Job can have many Tasks"),
		edge.To("tags", Tag.Type).
			Comment("A Job can have many Tags"),
		edge.To("next", Job.Type).
			Unique().
			From("prev").
			Unique(),
		edge.From("owner", User.Type).
			Ref("jobs").
			Unique().
			Required().
			Comment("A Job must have an owner"),
	}
}
