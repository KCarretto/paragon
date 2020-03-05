package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/kcarretto/paragon/ent"
)

// CreateUserSession creates a new user session for a request. It adds the authenticated user to the
// request context and ensures session cookies are set. If user is nil, this is a no-op.
func CreateUserSession(w http.ResponseWriter, req *http.Request, user *ent.User) *http.Request {
	if user == nil {
		return req
	}

	// Create session token for user if none exists
	if user.SessionToken == "" {
		user = user.Update().
			SetSessionToken(string(NewSecret(SessionTokenLength))).
			SaveX(req.Context())
	}

	// Ensure session token cookie is set
	if cookie, err := req.Cookie(SessionCookieName); err != nil || cookie == nil || cookie.Value == "" {
		http.SetCookie(w, &http.Cookie{
			Name:     SessionCookieName,
			Value:    user.SessionToken,
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Now().AddDate(0, 0, 1),
		})
	}

	// Ensure user id cookie is set
	if cookie, err := req.Cookie(UserCookieName); err != nil || cookie == nil || cookie.Value == "" {
		http.SetCookie(w, &http.Cookie{
			Name:    UserCookieName,
			Value:   fmt.Sprintf("%d", user.ID),
			Path:    "/",
			Expires: time.Now().AddDate(0, 0, 1),
		})
	}

	// Add user to request context
	return req.WithContext(context.WithValue(req.Context(), userContextKey, user))
}

// ClearUserSession clears any user session associated with the request. It is a no-op if no user
// session is associated with the request.
func ClearUserSession(w http.ResponseWriter) {
	// Expire session token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
	})

	// Expire user id cookie
	http.SetCookie(w, &http.Cookie{
		Name:    UserCookieName,
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
	})
}
