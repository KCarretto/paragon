package main

import (
	"context"
	"fmt"
	"os"

	gohttp "net/http"

	"github.com/kcarretto/paragon/graphql"
	"github.com/kcarretto/paragon/graphql/models"
	"github.com/kcarretto/paragon/pkg/agent/c2"
	"github.com/kcarretto/paragon/pkg/agent/transport/http"
	"github.com/kcarretto/paragon/pkg/agent/transport/websocket"
	"go.uber.org/zap"
)

// EnvDebug is the environment variable that determines if debug mode is enabled.
const (
	EnvDebug = "DEBUG"
)

func isDebug() bool {
	return os.Getenv(EnvDebug) != ""
}

func getLogger() (logger *zap.Logger) {
	var err error

	if isDebug() {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		panic(err)
	}

	return
}

func main() {
	logger := getLogger().Named("c2")

	// C2 Server Address
	httpAddr := "127.0.0.1:8080"
	if addr := os.Getenv("PG_HTTP_ADDR"); addr != "" {
		httpAddr = addr
	}

	// Teamserver URL
	teamserverHost := "127.0.0.1"
	teamserverURL := fmt.Sprintf("http://%s/graphql", teamserverHost)
	if url := os.Getenv("TEAMSERVER_URL"); url != "" {
		teamserverURL = url
	}
	graph := &graphql.Client{
		Service: "pg-c2",
		URL:     teamserverURL,
	}

	// Initial call to register the service
	graph.ClaimTasks(context.Background(), models.ClaimTasksRequest{
		PrimaryIP: nil,
	})

	// Initialize Server
	srv := &c2.Server{
		Log:        logger,
		Teamserver: graph,
	}

	mux := gohttp.NewServeMux()

	wsSvc := &websocket.ServerTransport{
		Log:        logger.Named("transport.wsgiveshell"),
		Handler:    websocket.HandleGiveShell,
		TeamServer: teamserverHost,
	}
	mux.Handle("/websocketgiveshell", wsSvc)

	httpSvc := &http.ServerTransport{
		Log:    logger.Named("transport.http"),
		Server: srv,
	}
	mux.Handle("/", httpSvc)

	if err := gohttp.ListenAndServe(httpAddr, mux); err != nil {
		panic(err)
	}
}
