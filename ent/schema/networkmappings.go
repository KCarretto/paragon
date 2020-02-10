package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
)

// NetworkMappings holds the schema definition for the NetworkMappings entity.
type NetworkMappings struct {
	ent.Schema
}

// Fields of the NetworkMappings.
func (NetworkMappings) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("Type").
			Values("LoadBalancer", "One To One").
			Comment("Type of mapping"),
	}
}

// Edges of the NetworkMappings.
func (NetworkMappings) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("Provider", Device.Type).
			Comment("The device providing the mapping"),
		// Potentially maps to devices too
	}
}
