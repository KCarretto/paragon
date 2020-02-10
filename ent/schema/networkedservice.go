package schema

import "github.com/facebookincubator/ent"
import "github.com/facebookincubator/ent/schema/field"

// NetworkedService holds the schema definition for the NetworkedService entity.
type NetworkedService struct {
	ent.Schema
}

// Fields of the NetworkedService.
func (NetworkedService) Fields() []ent.Field {
	return []ent.Field{
		field.String("Protocol"), // Enum tcp udp
		field.String("Name"),
		field.Int("Port"),
	}
}

// Edges of the NetworkedService.
func (NetworkedService) Edges() []ent.Edge {
	return nil
}
