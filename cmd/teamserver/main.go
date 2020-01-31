package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kcarretto/paragon/pkg/auth"
	"github.com/kcarretto/paragon/pkg/teamserver"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// HTTPAddr for the teamserver to listen on. Configure with HTTP_ADDR env variable.
var (
	start time.Time

	HTTPAddr          = "127.0.0.1:8080"
	OAuthClientID     string
	OAuthClientSecret auth.Secret
	OAuthDomain       = "http://127.0.0.1:8080"
)

// DefaultTopic is the default pubsub topic where events will be published. It can be configured
// by setting the PUB_TOPIC environment variable.
const DefaultTopic = "TOPIC"

func init() {
	start = time.Now()
	if addr := os.Getenv("PG_HTTP_ADDR"); addr != "" {
		HTTPAddr = addr
	}
	if clientID := os.Getenv("OAUTH_CLIENT_ID"); clientID != "" {
		OAuthClientID = clientID
	}
	if clientSecret := os.Getenv("OAUTH_CLIENT_SECRET"); clientSecret != "" {
		OAuthClientSecret = auth.Secret(clientSecret)
	}
	if domain := os.Getenv("OAUTH_DOMAIN"); domain != "" {
		OAuthDomain = domain
	}
}

// ServeHTTP handles http requests for the Teamserver.
func ServeHTTP(logger *zap.Logger, addr string, svc *teamserver.Service) error {
	router := http.NewServeMux()
	svc.HTTP(router)

	defer func() {
		logger.Info("Teamserver Stopped",
			zap.Duration("uptime", time.Since(start)),
		)
	}()

	logger.Info("Teamserver Started",
		zap.String("listen_on", addr),
		zap.Duration("start_latency", time.Since(start)),
	)

	return serve(addr, router)
}

func main() {
	ctx := context.Background()
	logger := newLogger().Named("svc.teamserver")
	logger.Debug("Initializing teamserver", zap.Time("start", start))

	topic := os.Getenv("PUB_TOPIC")
	if topic == "" {
		log.Println("[WARN] No PUB_TOPIC environment variable set to publish events, is this a mistake?")
	}
	publisher := newPublisher(ctx, topic)

	graph := newGraph(ctx, logger.Named("db"))
	defer graph.Close()

	authenticator := getAuthenticator(logger.Named("auth"), graph)

	svc := &teamserver.Service{
		Auth:   authenticator,
		Log:    logger,
		Graph:  graph,
		Events: publisher,
		OAuth: &oauth2.Config{
			ClientID:     OAuthClientID,
			ClientSecret: string(OAuthClientSecret),
			RedirectURL:  OAuthDomain + "/oauth/authorize",
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}

	if err := ServeHTTP(logger.Named("http"), HTTPAddr, svc); err != nil {
		logger.Fatal("HTTP server failure", zap.Error(err))
	}
}
