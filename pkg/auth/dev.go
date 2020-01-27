// +build dev

package auth

import (
	"net/http"

	"github.com/kcarretto/paragon/ent"
)

type DevAuthenticator struct {
	Graph *ent.Client
}

func (auth DevAuthenticator) Authenticate(req *http.Request) (*ent.User, error) {
	if userID, err := parseUserID(req); err == nil && userID > 0 {
		return auth.Graph.User.GetX(req.Context(), userID), nil
	}

	user, err := auth.Graph.User.Query().First(req.Context())
	if err != nil || user == nil {
		// TODO: Setup TOFU for admin
		// num := svc.Graph.User.Query().CountX(req.Context())

		user = auth.Graph.User.Create().
			SetName(string(NewSecret(10))).
			SetOAuthID(string(NewSecret(512))).
			SetPhotoURL("").
			SetSessionToken(string(NewSecret(SessionTokenLength))).
			SetIsAdmin(true).
			SetActivated(true).
			SaveX(req.Context())
		return user, nil
	}

	return user, nil
}
