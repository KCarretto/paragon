package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
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
			Comment("The ed25519 public key for the service (stored in Base64 of DER format)"),
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
