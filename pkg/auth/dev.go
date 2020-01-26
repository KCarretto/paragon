// +build dev

package auth

import (
	"fmt"
	"net/http"

	"github.com/kcarretto/paragon/ent"
)

type DevAuthenticator struct {
	Graph *ent.Client
}

func (auth DevAuthenticator) Authenticate(req *http.Request) (*ent.User, error) {
	var user *ent.User

	if userID, err := parseUserID(req); err == nil {
		fmt.Println("Parsed user ID", userID)
		user = auth.Graph.User.GetX(req.Context(), userID)
	} else {
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
		}
	}

	return user, nil
}
