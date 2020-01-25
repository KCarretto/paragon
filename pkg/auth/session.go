package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func (svc *Service) Authenticate(req *http.Request) (*Context, error) {
	userID, err := parseUserID(req)
	if err != nil {
		return nil, err
	}

	token, err := parseSessionToken(req)
	if err != nil {
		return nil, err
	}

	return svc.AuthenticateUser(req.Context(), userID, token)
}

func (svc *Service) AuthenticateUser(ctx context.Context, userID int, token Secret) (*Context, error) {
	user, err := svc.Graph.User.Get(ctx, userID)
	if err != nil || user == nil {
		return nil, fmt.Errorf("invalid user cookie: %w", err)
	}

	if user.SessionToken == "" {
		return nil, fmt.Errorf("user must reauthenticate")
	}

	if !token.Equals(Secret(user.SessionToken)) {
		return nil, fmt.Errorf("invalid session cookie")
	}

	return &Context{
		Context: ctx,
		user:    user,
	}, nil
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
	if userCookie.Expires.Before(time.Now()) {
		return -1, fmt.Errorf("user cookie expired")
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

	if sessionCookie.Expires.Before(time.Now()) {
		return "", fmt.Errorf("session cookie expired")
	}

	sessionToken := Secret(sessionCookie.Value)
	if sessionToken == "" {
		return "", fmt.Errorf("invalid session cookie")
	}

	return sessionToken, nil
}
