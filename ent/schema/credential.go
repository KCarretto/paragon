package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
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
		field.Text("secret").
			MaxLen(3000).
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
