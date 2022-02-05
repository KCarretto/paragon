package schema

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Link holds information about a link (usually to a file).
type Link struct {
	ent.Schema
}

// Fields of the Link.
func (Link) Fields() []ent.Field {
	return []ent.Field{
		field.String("Alias").
			NotEmpty().
			MinLen(1).
			Unique().
			Validate(func(val string) error {
				if strings.Contains(val, "/") {
					return fmt.Errorf("alias for link cannot contain slashes")
				}
				return nil
			}).
			Comment("The alias of the link which will be used for routing resolution"),
		field.Time("ExpirationTime").
			Optional().
			Comment("The timestamp for when the link will be deleted"),
		field.Int("Clicks").
			Default(-1).
			Min(-1).
			Comment("The number of clicks left on the link before it will be deleted (-1 means unlimited)"),
	}
}

// Edges of the Link.
func (Link) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("file", File.Type).
			Ref("links").
			Unique().
			Comment("A Link might have a file"),
	}
}
