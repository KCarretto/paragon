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

	user, err := auth.Graph.User.Get(req.Context(), userID)
	if err != nil || user == nil {
		return nil, fmt.Errorf("invalid user cookie: %w", err)
	}

	if token == "" || !token.Equals(Secret(user.SessionToken)) {
		return nil, fmt.Errorf("invalid session cookie")
	}

	if !user.Activated && !user.IsAdmin {
		return nil, fmt.Errorf("user pending activation")
	}

	return context.WithValue(req.Context(), userContextKey, user), nil
}

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
