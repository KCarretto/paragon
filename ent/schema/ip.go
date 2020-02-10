package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/field"
)

// IP holds the schema definition for the IP entity.
type IP struct {
	ent.Schema
}

// Fields of the IP.
func (IP) Fields() []ent.Field {
	return []ent.Field{
		field.String("IPAddress").
			Comment("The actual address"),
		field.Bool("Internal").
			Comment("Is this IP internal"), // what about internal to infra like container ips in kube
	}
}

// Edges of the IP.
func (IP) Edges() []ent.Edge {
	return nil
	// Overloaded to device for attacks (giving the router ip to a device)
	// IP can have ports open (services that have ports attached?)
}
