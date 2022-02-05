package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Tag holds the schema definition for the Tag entity.
type Tag struct {
	ent.Schema
}

// Fields of the Tag.
func (Tag) Fields() []ent.Field {
	return []ent.Field{
		field.String("Name").
			NotEmpty().
			Unique().
			Comment("The name of the Tag"),
	}
}

// Edges of the Tag.
func (Tag) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("targets", Target.Type).
			Ref("tags"),
		edge.From("tasks", Task.Type).
			Ref("tags"),
		edge.From("jobs", Job.Type).
			Ref("tags"),
	}
}
