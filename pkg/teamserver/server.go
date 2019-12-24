package teamserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/99designs/gqlgen/handler"
	"github.com/kcarretto/paragon/graphql/generated"
	"github.com/kcarretto/paragon/graphql/resolve"


	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/pkg/middleware"
)

// Server handles c2 messages and replies with new tasks for the c2 to send out.
type Server struct {
	Log       *zap.Logger
	EntClient *ent.Client
}

func (srv *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"status": "OK",
	}
	resp, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "failed to marshal the json for the status", http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(resp); err != nil {
		http.Error(w, fmt.Sprintf("Failed to write response data: %s", err.Error()), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
}

// Run begins the handlers for processing the subscriptions to the `tasks.claimed` and `tasks.executed` topics
func (srv *Server) Run() {
	router := http.NewServeMux()

	h := handler.GraphQL(generated.NewExecutableSchema(generated.Config{Resolvers: &resolve.Resolver{EntClient: srv.EntClient}}))

	router.Handle("/api/v1/", h)
	router.Handle("/api/playground/", handler.Playground("GraphQL", "/api/v1/"))
	router.HandleFunc("/status", srv.handleStatus)

	if err := http.ListenAndServe("0.0.0.0:80", middleware.Chain(router, middleware.WithPanicHandling)); err != nil {
		panic(err)
	}
}
