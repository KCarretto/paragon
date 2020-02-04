package teamserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/graphql"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"github.com/kcarretto/paragon/pkg/auth"
	"github.com/kcarretto/paragon/pkg/auth/oauth"
	"github.com/kcarretto/paragon/pkg/cdn"
	"github.com/kcarretto/paragon/pkg/event"
	"github.com/kcarretto/paragon/pkg/service"
	"github.com/kcarretto/paragon/www"
)

// Service provides HTTP handlers to compose the CDN, GraphQL, and WWW services.
type Service struct {
	Log    *zap.Logger
	Graph  *ent.Client
	Events event.Publisher
	OAuth  *oauth2.Config
	Auth   service.Authenticator
}

// HandleStatus returns JSON status: OK if the teamserver is running without error.
func (svc *Service) HandleStatus(w http.ResponseWriter, r *http.Request) error {
	type Response struct {
		Available       bool `json:"available"`
		UserID          int  `json:"userid,omitempty"`
		IsAuthenticated bool `json:"is_authenticated"`
		IsActivated     bool `json:"is_activated"`
		IsAdmin         bool `json:"is_admin"`
	}
	dst := json.NewEncoder(w)

	resp := Response{
		Available: true,
	}
	if user := auth.GetUser(r.Context()); user != nil {
		resp.UserID = user.ID
		resp.IsAuthenticated = true
		resp.IsActivated = user.IsActivated
		resp.IsAdmin = user.IsAdmin
	}

	if err := dst.Encode(resp); err != nil {
		return fmt.Errorf("failed to write response data: %w", err)
	}
	w.Header().Set("Content-Type", "application/json")
	return nil
}

// HTTP registers http handlers for the Teamserver.
func (svc *Service) HTTP(router *http.ServeMux) {
	oauthSVC := &oauth.Service{
		Log:    svc.Log.Named("auth"),
		Graph:  svc.Graph,
		Config: svc.OAuth,
	}
	graphqlSVC := &graphql.Service{
		Log:    svc.Log.Named("graphql"),
		Graph:  svc.Graph,
		Events: svc.Events,
		Auth:   svc.Auth,
	}
	cdnSVC := &cdn.Service{
		Log:   svc.Log.Named("cdn"),
		Graph: svc.Graph,
		Auth:  svc.Auth,
	}
	wwwSVC := &www.Service{
		Log:  svc.Log.Named("www"),
		Auth: svc.Auth,
	}
	status := &service.Endpoint{
		Log:           svc.Log.Named("status"),
		Authenticator: svc.Auth,
		Handler:       service.HandlerFn(svc.HandleStatus),
	}

	graphqlSVC.HTTP(router)
	cdnSVC.HTTP(router)
	wwwSVC.HTTP(router)

	router.Handle("/status", status)
	oauthSVC.HTTP(router)
}
