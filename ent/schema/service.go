package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Service holds the schema definition for the Service entity.
type Service struct {
	ent.Schema
}

// Fields of the Service.
func (Service) Fields() []ent.Field {
	return []ent.Field{
		field.String("Name").
			NotEmpty().
			Comment("The name displayed for the service"),
		field.String("PubKey").
			Unique().
			MaxLen(250).
			Comment("The ed25519 public key for the service (stored in Base64 of DER format)"),
		field.Text("Config").
			Default("").
			Comment("The configuration script of the service (usually a Renegade Script)"),
		field.Bool("IsActivated").
			Default(false).
			Comment("True iff the service is active and able to authenticate"),
	}
}

// Edges of the Service.
func (Service) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tag", Tag.Type).
			Unique().
			Required().
			Comment("A Service has a single tag"),
		edge.To("events", Event.Type).
			Comment("A Service can have many events"),
	}
}
