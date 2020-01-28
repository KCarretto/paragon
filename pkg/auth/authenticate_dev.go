// +build dev

package auth

import (
	"net/http"

	"github.com/kcarretto/paragon/ent"
)

// DevAuthenticator is only available when run in developer mode. It inserts users into the request
// context, even on unauthenticated requests. IT SHOULD NOT BE USED IN PRODUCTION.
type DevAuthenticator struct {
	Graph *ent.Client
}

// Authenticate in developer mode, which creates users for unauthenticated requests. Only returns
// transient errors.
func (auth DevAuthenticator) Authenticate(w http.ResponseWriter, req *http.Request) (*http.Request, error) {
	// If a userID was provided, attempt to create a session for that user
	if userID, err := parseUserID(req); err == nil && userID > 0 {
		user, err := auth.Graph.User.Get(req.Context(), userID)
		if err == nil && user != nil {
			return CreateUserSession(w, req, user), nil
		}
	}

	// Invalid or no user id provided, assign the session to the first user
	user, err := auth.Graph.User.Query().First(req.Context())
	if err == nil && user != nil {
		return CreateUserSession(w, req, user), nil
	}

	// Otherwise, create a user for the session
	user = auth.Graph.User.Create().
		SetName(string(NewSecret(10))).
		SetOAuthID(string(NewSecret(512))).
		SetPhotoURL("").
		SetSessionToken(string(NewSecret(SessionTokenLength))).
		SetIsAdmin(true).
		SetIsActivated(true).
		SaveX(req.Context())
	return CreateUserSession(w, req, user), nil
}
