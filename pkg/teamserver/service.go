package teamserver

import (
	"fmt"
	"io"
	"net/http"

	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/graphql"
	"go.uber.org/zap"

	"github.com/kcarretto/paragon/pkg/cdn"
	"github.com/kcarretto/paragon/pkg/event"
	"github.com/kcarretto/paragon/www"
)

// Service provides HTTP handlers to compose the CDN, GraphQL, and WWW services.
type Service struct {
	Log    *zap.Logger
	Graph  *ent.Client
	Events event.Publisher
}

// HandleStatus returns JSON status: OK if the teamserver is running without error.
func (svc *Service) HandleStatus(w http.ResponseWriter, r *http.Request) {
	if _, err := io.WriteString(w, `{"status":"OK"}`); err != nil {
		http.Error(
			w,
			fmt.Sprintf("failed to write response data: %s", err.Error()),
			http.StatusInternalServerError,
		)
	}
	w.Header().Set("Content-Type", "application/json")
}

// HTTP registers http handlers for the Teamserver.
func (svc *Service) HTTP(router *http.ServeMux) {
	graphqlSVC := &graphql.Service{
		Log:    svc.Log.Named("graphql"),
		Graph:  svc.Graph,
		Events: svc.Events,
	}
	cdnSVC := &cdn.Service{
		Log:   svc.Log.Named("cdn"),
		Graph: svc.Graph,
	}
	wwwSVC := &www.Service{
		Log: svc.Log.Named("www"),
	}

	router.HandleFunc("/status", svc.HandleStatus)

	graphqlSVC.HTTP(router)
	cdnSVC.HTTP(router)
	wwwSVC.HTTP(router)
}
