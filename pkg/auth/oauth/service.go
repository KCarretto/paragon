package oauth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kcarretto/paragon/ent"
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
	exists, err := userQuery.Exist(req.Context())
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
			SetEmail(profile.Email).
			SetPhotoURL(profile.PhotoURL).
			SetIsAdmin(num == 0).
			SetActivated(num == 0).
			SaveX(req.Context())
	}

	if usr.SessionToken == "" {
		usr = usr.Update().
			SetSessionToken(string(auth.NewSecret(auth.SessionTokenLength))).
			SaveX(req.Context())
	}

	http.SetCookie(w, &http.Cookie{
		Name:     auth.SessionCookieName,
		Value:    usr.SessionToken,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().AddDate(0, 0, 1),
	})
	http.SetCookie(w, &http.Cookie{
		Name:    auth.UserCookieName,
		Value:   fmt.Sprintf("%d", usr.ID),
		Path:    "/",
		Expires: time.Now().AddDate(0, 0, 1),
	})
	http.Redirect(w, req, "/index.html", http.StatusTemporaryRedirect)
	return nil
}

// type SignupRequest struct {
// 	Name string `json:"username"`
// }

// type SignupResponse struct {
// 	Enabled bool   `json:"enabled"`
// 	URL     string `json:"url"`
// }

// type OAuthServer struct {
// 	ClientID     string `json:"client_id"`
// 	ClientSecret Secret `json:"client_secret"`
// 	RedirectURL  string `json:"-"`

// 	Graph *ent.Client `json:"-"`

// 	oauthConfig *oauth2.Config
// }

// func (srv *OAuthServer) state() string {
// 	b := make([]byte, OAuthStateLength)
// 	if _, err := rand.Read(b); err != nil {
// 		panic(err)
// 	}
// 	return base64.StdEncoding.EncodeToString(b)
// }

// func (srv *OAuthServer) sessionToken() Secret {
// b := make([]byte, SessionTokenLength)
// if _, err := rand.Read(b); err != nil {
// 	panic(err)
// }
// return Secret(base64.StdEncoding.EncodeToString(b))
// }

// func (srv *OAuthServer) config() *oauth2.Config {
// 	if srv.oauthConfig == nil {
// 		srv.oauthConfig = &oauth2.Config{
// 			ClientID:     srv.ClientID,
// 			ClientSecret: string(srv.ClientSecret),
// 			Endpoint:     google.Endpoint,
// 			RedirectURL:  srv.RedirectURL,
// 			Scopes: []string{
// 				"https://www.googleapis.com/auth/userinfo.email", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
// 			},
// 		}
// 	}
// 	return srv.oauthConfig
// }

// func (srv *OAuthServer) HandleSignup(w http.ResponseWriter, req *http.Request) {
// 	encoder := json.NewEncoder(w)
// 	decoder := json.NewDecoder(req.Body)

// 	info := SignupRequest{}
// 	if err := decoder.Decode(&info); err != nil {
// 		http.Error(w, "failed to decode json request", http.StatusBadRequest)
// 		return
// 	}

// 	state := srv.state()
// 	url := srv.config().AuthCodeURL(state)

// 	usr := srv.Graph.User.Create().
// 		SetName(info.Name).
// 		SetOAuthState(state).
// 		SaveX(req.Context())

// 	if err := encoder.Encode(&SignupResponse{
// 		Enabled: true,
// 		URL:     url,
// 	}); err != nil {
// 		http.Error(w, "failed to encode json response", http.StatusInternalServerError)
// 		srv.Graph.User.DeleteOne(usr).ExecX(req.Context())
// 	}
// }

// func (srv *OAuthServer) HandleAuth(w http.ResponseWriter, req *http.Request) {
// 	state := req.URL.Query().Get("state")
// 	id, err := srv.Graph.User.Query().Where(user.OAuthState(state)).OnlyID(req.Context())
// 	if err != nil {
// 		http.Error(w, "invalid oauth2 state", http.StatusBadRequest)
// 		fmt.Printf("\n\n\nOAUTH USER LOOKUP ERROR: %+v\n\n\n", err)
// 		return
// 	}

// 	code := req.URL.Query().Get("code")
// 	token, err := srv.config().Exchange(req.Context(), code)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("failed to validate oauth2 code: %v", err), http.StatusBadRequest)
// 		return
// 	}

// 	client := srv.config().Client(req.Context(), token)
// 	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("failed to retrieve oauth profile info: %v", err), http.StatusBadRequest)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	decoder := json.NewDecoder(resp.Body)
// 	data := userInfo{}

// 	if err := decoder.Decode(&data); err != nil {
// 		http.Error(w, fmt.Sprintf("failed to parse oauth profile info: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// session := srv.sessionToken()
// usr := srv.Graph.User.UpdateOneID(id).
// 	SetSessionToken(string(session)).
// 	SetActivated(true).
// 	SaveX(req.Context())

// http.SetCookie(w, &http.Cookie{
// 	Name:     SessionCookieName,
// 	Value:    string(session),
// 	HttpOnly: true,
// })
// http.SetCookie(w, &http.Cookie{
// 	Name:  UserCookieName,
// 	Value: fmt.Sprintf("%d", usr.ID),
// })
// http.Redirect(w, req, "/index.html", http.StatusTemporaryRedirect)
// }

// func NewOAuthServer(graph *ent.Client) *OAuthServer {
// 	id := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
// 	secret := Secret(os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"))
// 	if id == "" || secret == "" {
// 		panic("Must set GOOGLE_OAUTH_CLIENT_ID and GOOGLE_OAUTH_CLIENT_SECRET!")
// 	}

// 	return &OAuthServer{
// 		ClientID:     id,
// 		ClientSecret: secret,
// 		RedirectURL:  "http://localhost:80/oauth/authorize",
// 		Graph:        graph,
// 	}
// }
