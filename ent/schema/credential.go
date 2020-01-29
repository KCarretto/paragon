package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
)

// Credential holds the schema definition for the Credential entity.
type Credential struct {
	ent.Schema
}

// Fields of the Credential.
func (Credential) Fields() []ent.Field {
	return []ent.Field{
		field.String("principal").
			Comment("The principal for the Credential"),
		field.String("secret").
			Comment("The secret for the Credential"),
		field.Enum("kind").
			Values("password", "key", "certificate").
			Comment("The kind of the credential (password, key, etc)"),
		field.Int("fails").
			Default(0).
			Min(0).
			Comment("The number of failures for the Credential"),
	}
}

// Edges of the Credential.
func (Credential) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("target", Target.Type).
			Ref("credentials").
			Unique(),
	}
}
