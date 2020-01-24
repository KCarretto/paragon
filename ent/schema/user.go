package schema

import (
	"github.com/facebookincubator/ent"
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
			Unique().
			Comment("The name displayed for the user"),
		field.String("OAuthState").
			Sensitive().
			Unique().
			Comment("State used during OAuth signup"),
		field.String("SessionToken").
			Optional().
			Sensitive().
			Comment("The session token currently authenticating the user"),
		field.Bool("Activated").
			Default(false).
			Comment("True if the user is active and able to authenticate"),
	}
}
