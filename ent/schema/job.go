package schema

import (
	"time"

	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
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
		field.String("Content").
			NotEmpty().
			Comment("The content of the job (usually a Renegade Script)"),
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
