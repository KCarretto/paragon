package schema

import (
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
		field.String("Parameters").
			Comment("The JSON string for the Parameters Mapping"),
	}
}

// Edges of the Job.
func (Job) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tasks", Task.Type).
			Comment("A Job can have many Tasks"),
		edge.To("tags", Tag.Type).
			Comment("A Job can have many Tags"),
		edge.From("template", JobTemplate.Type).
			Ref("jobs").
			Unique().
			Required().
			Comment("A Job must have one template"),
	}
}
