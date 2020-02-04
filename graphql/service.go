package graphql

import (
	"net/http"

	"github.com/99designs/gqlgen/handler"
	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/graphql/generated"
	"github.com/kcarretto/paragon/graphql/resolve"
	"github.com/kcarretto/paragon/pkg/auth"
	"github.com/kcarretto/paragon/pkg/event"
	"github.com/kcarretto/paragon/pkg/service"

	"go.uber.org/zap"
)

// Service provides HTTP handlers for the GraphQL schema.
type Service struct {
	Log    *zap.Logger
	Graph  *ent.Client
	Events event.Publisher
	Auth   service.Authenticator
}

// HandleGraphQL initializes and returns a new GraphQL API handler.
func (svc *Service) HandleGraphQL() http.HandlerFunc {
	resolver := &resolve.Resolver{
		Log:    svc.Log.Named("resolver"),
		Graph:  svc.Graph,
		Events: svc.Events,
	}
	config := generated.Config{Resolvers: resolver}
	schema := generated.NewExecutableSchema(config)

	return handler.GraphQL(schema)
}

// HandlePlayground initializes and returns a new GraphQL Playground handler.
func (svc *Service) HandlePlayground() http.HandlerFunc {
	return handler.Playground("GraphQL", "/graphql")
}

// HTTP registers http handlers for a GraphQL API.
func (svc *Service) HTTP(router *http.ServeMux) {
	api := &service.Endpoint{
		Log:           svc.Log.Named("api"),
		Authenticator: svc.Auth,
		Authorizer:    auth.NewAuthorizer().IsActivated(),
		Handler:       service.HTTPHandler(svc.HandleGraphQL()),
	}
	graphiql := &service.Endpoint{
		Log:           svc.Log.Named("graphiql"),
		Authenticator: svc.Auth,
		Authorizer:    auth.NewAuthorizer().IsActivated(),
		Handler:       service.HTTPHandler(svc.HandlePlayground()),
	}
	router.Handle("/graphql", api)
	router.Handle("/graphiql", graphiql)
}
