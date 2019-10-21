package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
)

// Target holds the schema definition for the Target entity.
type Target struct {
	ent.Schema
}

// Fields of the Target.
func (Target) Fields() []ent.Field {
	return []ent.Field{
		field.String("Name").
			Optional().
			Comment("The name of the Target"),
		field.String("MachineUUID").
			Unique().
			Comment("The machine UUID of the Target"),
		field.String("PrimaryIP").
			Optional().
			Comment("The IP Address for the primary interface of the Target"),
		field.String("PublicIP").
			Optional().
			Comment("The Public IP Address for the Target"),
		field.String("PrimaryMAC").
			Optional().
			Comment("The MAC Address for the primary interface of the Target"),
		field.String("Hostname").
			Optional().
			Comment("The hostname for the Target"),
		field.Time("LastSeen").
			Optional().
			Comment("The time the Target was last seen"),
	}
}

// Edges of the Target.
func (Target) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tasks", Task.Type).
			Comment("A Target can have many tasks connected to it"),
	}
}
