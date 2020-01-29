package auth

import (
	"context"

	"github.com/kcarretto/paragon/ent"
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
	svc := GetService(ctx)

	if svc != nil {
		return authz.authorizeService(svc)
	}

	if user != nil {
		return authz.authorizeUser(user)
	}

	return ErrNotAuthenticated
}

func (authz *Authorizer) authorizeUser(user *ent.User) error {
	if authz.requireActivated && !user.IsActivated {
		return ErrUnauthorized
	}

	if authz.requireAdmin && !user.IsAdmin {
		return ErrUnauthorized
	}
	return nil
}

func (authz *Authorizer) authorizeService(svc *ent.Service) error {
	if authz.requireActivated && !svc.IsActivated {
		return ErrUnauthorized
	}
	return nil
}
