package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/ent/user"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type userInfo struct {
	Email      string `json:"email"`
	Name       string `json:"name"`
	IsVerified bool   `json:"email_verified"`
}

type SignupRequest struct {
	Name string `json:"username"`
}

type SignupResponse struct {
	Enabled bool   `json:"enabled"`
	URL     string `json:"url"`
}

type OAuthServer struct {
	ClientID     string `json:"client_id"`
	ClientSecret Secret `json:"client_secret"`
	RedirectURL  string `json:"-"`

	Graph *ent.Client `json:"-"`

	oauthConfig *oauth2.Config
}

func (srv *OAuthServer) state() string {
	b := make([]byte, OAuthStateLength)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}

func (srv *OAuthServer) sessionToken() Secret {
	b := make([]byte, SessionTokenLength)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return Secret(base64.StdEncoding.EncodeToString(b))
}

func (srv *OAuthServer) config() *oauth2.Config {
	if srv.oauthConfig == nil {
		srv.oauthConfig = &oauth2.Config{
			ClientID:     srv.ClientID,
			ClientSecret: string(srv.ClientSecret),
			Endpoint:     google.Endpoint,
			RedirectURL:  srv.RedirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
			},
		}
	}
	return srv.oauthConfig
}

func (srv *OAuthServer) HandleSignup(w http.ResponseWriter, req *http.Request) {
	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(req.Body)

	info := SignupRequest{}
	if err := decoder.Decode(&info); err != nil {
		http.Error(w, "failed to decode json request", http.StatusBadRequest)
		return
	}

	state := srv.state()
	url := srv.config().AuthCodeURL(state)

	usr := srv.Graph.User.Create().
		SetName(info.Name).
		SetOAuthState(state).
		SaveX(req.Context())

	if err := encoder.Encode(&SignupResponse{
		Enabled: true,
		URL:     url,
	}); err != nil {
		http.Error(w, "failed to encode json response", http.StatusInternalServerError)
		srv.Graph.User.DeleteOne(usr).ExecX(req.Context())
	}
}

func (srv *OAuthServer) HandleAuth(w http.ResponseWriter, req *http.Request) {
	state := req.URL.Query().Get("state")
	id, err := srv.Graph.User.Query().Where(user.OAuthState(state)).OnlyID(req.Context())
	if err != nil {
		http.Error(w, "invalid oauth2 state", http.StatusBadRequest)
		fmt.Printf("\n\n\nOAUTH USER LOOKUP ERROR: %+v\n\n\n", err)
		return
	}

	code := req.URL.Query().Get("code")
	token, err := srv.config().Exchange(req.Context(), code)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to validate oauth2 code: %v", err), http.StatusBadRequest)
		return
	}

	client := srv.config().Client(req.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to retrieve oauth profile info: %v", err), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	data := userInfo{}

	if err := decoder.Decode(&data); err != nil {
		http.Error(w, fmt.Sprintf("failed to parse oauth profile info: %v", err), http.StatusInternalServerError)
		return
	}

	session := srv.sessionToken()
	usr := srv.Graph.User.UpdateOneID(id).
		SetSessionToken(string(session)).
		SetActivated(true).
		SaveX(req.Context())

	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    string(session),
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:  UserCookieName,
		Value: fmt.Sprintf("%d", usr.ID),
	})
	http.Redirect(w, req, "/index.html", http.StatusTemporaryRedirect)
}

func NewOAuthServer(graph *ent.Client) *OAuthServer {
	id := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	secret := Secret(os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"))
	if id == "" || secret == "" {
		panic("Must set GOOGLE_OAUTH_CLIENT_ID and GOOGLE_OAUTH_CLIENT_SECRET!")
	}

	return &OAuthServer{
		ClientID:     id,
		ClientSecret: secret,
		RedirectURL:  "http://localhost:80/oauth/authorize",
		Graph:        graph,
	}
}
