//go:build dev
// +build dev

package auth

import (
	"context"
	"net/http"

	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/ent/event"
	"github.com/kcarretto/paragon/ent/service"
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

	// If a service was provided, upsert the service and allow it.
	if svcName := req.Header.Get(HeaderService); svcName != "" {
		svc, err := auth.Graph.Service.Query().
			Where(service.Name(svcName)).
			First(req.Context())

		// Upsert the service if not found
		if err != nil || svc == nil {
			tag, err := auth.Graph.Tag.Create().
				SetName(svcTagPrefix + svcName).
				Save(req.Context())
			if err != nil {
				return nil, err
			}

			svc, err = auth.Graph.Service.Create().
				SetName(svcName).
				SetPubKey(svcName).
				SetTag(tag).
				SetIsActivated(true).
				Save(req.Context())
			if err != nil {
				return nil, err
			}
		}
		return req.WithContext(context.WithValue(req.Context(), serviceContextKey, svc)), nil
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

	// we silent fail because as cool as they are, events are less important than functionality
	auth.Graph.Event.Create().
		SetOwner(user).
		SetUser(user).
		SetKind(event.KindCREATE_USER).
		Save(req.Context())
	return CreateUserSession(w, req, user), nil
}
