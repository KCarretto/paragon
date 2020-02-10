package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/field"
)

// Subnet holds the schema definition for the Subnet entity.
type Subnet struct {
	ent.Schema
}

// Fields of the Subnet.
func (Subnet) Fields() []ent.Field {
	return []ent.Field{
		field.String("Name").
			Comment("A name for the subnet"),
		field.String("NetworkAddress").
			Comment("The Network Address for the subnet"),
		field.String("Mask").
			Comment("The number of network bits of the subnet"),
	}
}

// Edges of the Subnet.
func (Subnet) Edges() []ent.Edge {
	return nil
}
