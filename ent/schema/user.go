package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("Name").
			MinLen(3).
			MaxLen(25).
			Comment("The name displayed for the user"),
		field.String("OAuthID").
			Sensitive().
			Unique().
			Immutable().
			Comment("OAuth Subject ID of the user"),
		field.String("PhotoURL").
			Comment("URL to the user's profile photo."),
		field.String("SessionToken").
			Optional().
			Sensitive().
			Comment("The session token currently authenticating the user"),
		field.Bool("IsActivated").
			Default(false).
			Comment("True iff the user is active and able to authenticate"),
		field.Bool("IsAdmin").
			Default(false).
			Comment("True iff the user is an Admin"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("jobs", Job.Type).
			Comment("A User can have many Jobs Created"),
		edge.To("events", Event.Type).
			Comment("A User can have many events"),
	}
}
