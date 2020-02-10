package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/field"
)

// Interface holds the schema definition for the Interface entity.
type Interface struct {
	ent.Schema
}

// Fields of the Interface.
func (Interface) Fields() []ent.Field {
	return []ent.Field{
		field.String("Name").
			Comment("The name of the Interface"),
		field.String("MAC").
			Optional().
			Comment("The MAC Address for the Interface"),
		field.Strings("IPAddresses").
			Optional().
			Comment("The ip addresses on the Interface"),
	}
}

// Edges of the Interface.
func (Interface) Edges() []ent.Edge {
	return nil
	// edge to ips instead
}
