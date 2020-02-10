package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
)

// Device holds the schema definition for the Device entity.
type Device struct {
	ent.Schema
}

// Fields of the Device.
func (Device) Fields() []ent.Field {
	return []ent.Field{
		field.String("Name").
			Comment("The name of the Device"),
		field.String("MachineUUID").
			MaxLen(250).
			Unique().
			Optional().
			Comment("The machine UUID of the Device"),
		field.String("Hostname").
			Optional().
			Comment("The hostname for the Device"),
		field.String("OS").
			Optional().
			Comment("The base OS of the device"),
		field.Time("LastUpdated").
			Optional().
			Comment("The time the Device was last updated"),
	}
}

// Edges of the Device.
func (Device) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("interfaces", Interface.Type).
			Comment("A Device can have multiple interfaces"),
		edge.To("PublicIPs", IP.Type).
			Comment("A Device can be linked to multiple public IPs"),
	}
}
