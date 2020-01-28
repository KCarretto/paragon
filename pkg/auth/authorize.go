package auth

import (
	"context"
)

// An Authorizer asserts various properties of a requesting context.
type Authorizer struct {
	requireActivated bool
	requireAdmin     bool
}

// NewAuthorizer initializes and returns a new authorizer.
func NewAuthorizer() *Authorizer {
	return &Authorizer{}
}

// IsAdmin ensures that authorized users have admin privileges.
func (authz *Authorizer) IsAdmin() *Authorizer {
	authz.requireAdmin = true
	return authz
}

// IsActivated ensures that authorized users are activated.
func (authz *Authorizer) IsActivated() *Authorizer {
	authz.requireActivated = true
	return authz
}

// Authorize the provided context based on the preconfigured rules.
func (authz *Authorizer) Authorize(ctx context.Context) error {
	user := GetUser(ctx)
	if user == nil {
		return ErrNotAuthenticated
	}

	if authz.requireActivated && !user.Activated {
		return ErrUnauthorized
	}

	if authz.requireAdmin && !user.Activated {
		return ErrUnauthorized
	}

	return nil
}
