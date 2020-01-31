package main

import (
	"net/http"
	"os"

	"github.com/kcarretto/paragon/graphql"
	"github.com/kcarretto/paragon/pkg/c2"

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
	logger := getLogger()

	// C2 Server Address
	httpAddr := "127.0.0.1:8080"
	if addr := os.Getenv("PG_HTTP_ADDR"); addr != "" {
		httpAddr = addr
	}

	// Teamserver URL
	teamserverURL := "http://127.0.0.1/graphql"
	if url := os.Getenv("TEAMSERVER_URL"); url != "" {
		teamserverURL = url
	}

	// Initialize Server
	srv := &c2.Server{
		Teamserver: &graphql.Client{
			Service: "pg-c2",
			URL:     teamserverURL,
		},
	}

	// Setup routes
	router := http.NewServeMux()
	router.Handle("/", srv)
	router.HandleFunc("/status", srv.ServeStatus)

	handler := withLogging(logger.Named("http"), router)

	logger.Info("Starting HTTP Server", zap.String("addr", httpAddr))
	if err := http.ListenAndServe(httpAddr, handler); err != nil {
		panic(err)
	}
}
