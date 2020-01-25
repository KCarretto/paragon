package auth

import (
	"context"

	"github.com/kcarretto/paragon/ent"
)

const (
	SessionTokenLength = 256
	SessionCookieName  = "pg-session"
	UserCookieName     = "pg-userid"
	LoginURL           = "/oauth/login"
)

// UserContext of an authenticated user.
type UserContext interface {
	UserID() int
	User() *ent.User
}

// Context of an authenticated identity.
type Context struct {
	context.Context

	user *ent.User
}

func (ctx *Context) UserID() int {
	return ctx.user.ID
}

func (ctx *Context) User() *ent.User {
	return ctx.user
}
