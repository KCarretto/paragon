package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
)

// JobTemplate holds the schema definition for the JobTemplate entity.
type JobTemplate struct {
	ent.Schema
}

// Fields of the Target.
func (JobTemplate) Fields() []ent.Field {
	return []ent.Field{
		field.String("Name").
			NotEmpty().
			Comment("The name of the job template"),
		field.String("Content").
			NotEmpty().
			Comment("The content of the job template (usually a Renegade Script)"),
	}
}

// Edges of the Job.
func (JobTemplate) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("jobs", Job.Type).
			Comment("A Job Template can have many Jobs"),
		edge.To("tags", Tag.Type).
			Comment("A Job Template can have many Tags"),
	}
}
