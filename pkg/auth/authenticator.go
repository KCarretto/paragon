package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kcarretto/paragon/ent"
)

const (
	SessionTokenLength = 256
	SessionCookieName  = "pg-session"
	UserCookieName     = "pg-userid"
)

type authKey string

var userContextKey authKey = "user"

type HTTPAuthenticator struct {
	Graph *ent.Client
}

func (auth HTTPAuthenticator) Authenticate(req *http.Request) (context.Context, error) {
	userID, err := parseUserID(req)
	if err != nil {
		return nil, err
	}

	token, err := parseSessionToken(req)
	if err != nil {
		return nil, err
	}

	return auth.AuthenticateUser(req.Context(), userID, token)
}

func (auth HTTPAuthenticator) AuthenticateUser(ctx context.Context, userID int, token Secret) (context.Context, error) {
	user, err := auth.Graph.User.Get(ctx, userID)
	if err != nil || user == nil {
		return nil, fmt.Errorf("invalid user cookie: %w", err)
	}

	if user.SessionToken == "" {
		return nil, fmt.Errorf("user must reauthenticate")
	}

	if !token.Equals(Secret(user.SessionToken)) {
		return nil, fmt.Errorf("invalid session cookie")
	}

	return context.WithValue(ctx, userContextKey, user), nil
}

func GetUser(ctx context.Context) *ent.User {
	if v := ctx.Value(userContextKey); v != nil {
		if usr, ok := v.(*ent.User); ok {
			return usr
		}
		panic(fmt.Errorf("Received non-user value for user context key: %v", v))
	}

	return nil
}

// func CreateUserSession(ctx context.Context, user *ent.User) *Context {
// 	if user.SessionToken == "" {
// 		token := NewSecret(SessionTokenLength)
// 		user = user.Update().SetSessionToken(string(token)).SaveX(ctx)
// 	}
// 	return &Context{
// 		Context:      ctx,
// 		userID:       user.ID,
// 		user:         user,
// 		sessionToken: Secret(user.SessionToken),
// 	}
// }

func parseUserID(req *http.Request) (int, error) {
	userCookie, err := req.Cookie(UserCookieName)
	if err != nil {
		return -1, fmt.Errorf("no user cookie set: %w", err)
	}

	userID, err := strconv.Atoi(userCookie.Value)
	if err != nil {
		return -1, fmt.Errorf("invalid user cookie: %w", err)
	}
	return userID, nil
}

func parseSessionToken(req *http.Request) (Secret, error) {
	sessionCookie, err := req.Cookie(SessionCookieName)
	if err != nil {
		return "", fmt.Errorf("no session cookie set: %w", err)
	}

	sessionToken := Secret(sessionCookie.Value)
	if sessionToken == "" {
		return "", fmt.Errorf("invalid session cookie")
	}

	return sessionToken, nil
}
