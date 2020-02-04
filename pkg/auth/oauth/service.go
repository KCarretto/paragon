package oauth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/ent/event"
	"github.com/kcarretto/paragon/ent/user"
	"github.com/kcarretto/paragon/pkg/auth"
	"github.com/kcarretto/paragon/pkg/service"
	"go.uber.org/zap"

	"golang.org/x/oauth2"
)

type userInfo struct {
	OAuthID    string `json:"sub"`
	PhotoURL   string `json:"picture"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	IsVerified bool   `json:"email_verified"`
}

const OAuthStateLength = 64

var oauth2State string

func init() {
	b := make([]byte, OAuthStateLength)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	oauth2State = base64.StdEncoding.EncodeToString(b)
}

// Service provides HTTP handlers for the Auth service.
type Service struct {
	Log    *zap.Logger
	Graph  *ent.Client
	Config *oauth2.Config
}

// HTTP registers http handlers for the Auth service.
func (svc *Service) HTTP(router *http.ServeMux) {
	login := &service.Endpoint{
		Log:     svc.Log.Named("login"),
		Handler: service.HandlerFn(svc.HandleLogin),
	}
	authorize := &service.Endpoint{
		Log:     svc.Log.Named("authorize"),
		Handler: service.HandlerFn(svc.HandleOAuth),
	}

	router.Handle("/oauth/login", login)
	router.Handle("/oauth/authorize", authorize)
}

// HandleLogin creates an OAuth 2.0 code url and redirects the client.
func (svc Service) HandleLogin(w http.ResponseWriter, req *http.Request) error {
	url := svc.Config.AuthCodeURL(oauth2State)
	http.Redirect(w, req, url, http.StatusTemporaryRedirect)
	return nil
}

// HandleOAuth authorizes users after being redirected from the OAuth consent screen with an access code.
func (svc Service) HandleOAuth(w http.ResponseWriter, req *http.Request) error {
	if state := req.URL.Query().Get("state"); state != oauth2State {
		return fmt.Errorf("invalid oauth2 state")
	}

	code := req.URL.Query().Get("code")
	token, err := svc.Config.Exchange(req.Context(), code)
	if err != nil {
		return fmt.Errorf("failed to validate oauth2 code: %w", err)
	}

	client := svc.Config.Client(req.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return fmt.Errorf("failed to retrieve oauth2 profile info: %w", err)
	}
	defer resp.Body.Close()

	var profile userInfo
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&profile); err != nil {
		return fmt.Errorf("failed to parse oauth2 profile info: %w", err)
	}
	if profile.OAuthID == "" {
		return fmt.Errorf("invalid oauth2 subject id: %w", err)
	}

	userQuery := svc.Graph.User.Query().Where(user.OAuthID(profile.OAuthID))
	exists, err := userQuery.Clone().Exist(req.Context())
	if err != nil {
		return fmt.Errorf("failed to retrieve user info: %w", err)
	}

	var usr *ent.User

	if exists {
		usr = userQuery.OnlyX(req.Context())
	} else {
		num := svc.Graph.User.Query().CountX(req.Context())

		usr = svc.Graph.User.Create().
			SetOAuthID(profile.OAuthID).
			SetSessionToken(string(auth.NewSecret(auth.SessionTokenLength))).
			SetName(profile.Name).
			SetPhotoURL(profile.PhotoURL).
			SetIsAdmin(num == 0).
			SetIsActivated(num == 0).
			SaveX(req.Context())

		// we silent fail because as cool as they are, events are less important than functionality
		svc.Graph.Event.Create().
			SetOwner(usr).
			SetUser(usr).
			SetKind(event.KindCREATEUSER).
			Save(req.Context())
	}

	req = auth.CreateUserSession(w, req, usr)
	http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
	return nil
}
