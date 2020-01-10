package teamserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"os"

	"go.uber.org/zap"

	"github.com/99designs/gqlgen/handler"
	"github.com/kcarretto/paragon/graphql/generated"
	"github.com/kcarretto/paragon/graphql/resolve"

	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/pkg/middleware"
	"github.com/kcarretto/paragon/www"
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

func (srv *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	f, err := www.App.Open("index.html")
	if err != nil {
		http.Error(w, "Failed to load index.html", 404)
		return
	}

	var modtime time.Time
	if info, err := f.Stat(); err == nil && info != nil {
		modtime = info.ModTime()
	}

	http.ServeContent(w, r, "index.html", modtime, f)
}

// Run begins the handlers for processing the subscriptions to the `tasks.claimed` and `tasks.executed` topics
func (srv *Server) Run() {
	router := http.NewServeMux()

	h := handler.GraphQL(generated.NewExecutableSchema(generated.Config{Resolvers: &resolve.Resolver{EntClient: srv.EntClient}}))

	router.Handle("/graphql", h)
	router.Handle("/graphiql", handler.Playground("GraphQL", "/graphql"))
	router.HandleFunc("/status", srv.handleStatus)
	router.Handle("/app/", http.StripPrefix("/app", http.FileServer(www.App)))
	router.HandleFunc("/", srv.handleIndex)

	// C2 Server Address
	httpAddr := "127.0.0.1:8080"
	if addr := os.Getenv("HTTP_ADDR"); addr != "" {
		httpAddr = addr
	}

	srv.Log.Info("Teamserver Initialized", zap.String("listen_on", httpAddr))
	if err := http.ListenAndServe(
		httpAddr,
		middleware.Chain(
			router,
			middleware.WithPanicHandling,
			middleware.WithLogging(srv.Log.Named("info")),
		)); err != nil {
		panic(err)
	}
}
