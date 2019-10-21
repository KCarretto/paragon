package main

import (
	"context"
	"net/http"
	"os"

	"github.com/kcarretto/paragon/c2"

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
	ctx := context.Background()
	logger := getLogger()

	execTopic, err := openTopic(ctx, "tasks.executed")
	if err != nil {
		logger.Panic("Failed to open pubsub topic", zap.Error(err))
	}
	defer execTopic.Shutdown(ctx)

	claimTopic, err := openTopic(ctx, "tasks.claimed")
	if err != nil {
		logger.Panic("Failed to open pubsub topic", zap.Error(err))
	}
	defer claimTopic.Shutdown(ctx)

	httpAddr := "127.0.0.1:8080"
	if addr := os.Getenv("HTTP_ADDR"); addr != "" {
		httpAddr = addr
	}

	srv := &c2.Server{
		OnTaskExecuted: onExec(ctx, logger.Named("events.tasks.executed"), execTopic),
		Queue: &c2.Queue{
			OnClaim: onClaim(ctx, logger.Named("events.tasks.claimed"), claimTopic),
		},
	}
	router := http.NewServeMux()
	router.Handle("/", srv)
	router.HandleFunc("/status", srv.ServeStatus)
	router.HandleFunc("/events/tasks/claimed", srv.ServeEventTaskClaimed)
	router.HandleFunc("/events/tasks/queued", srv.ServeEventTaskQueued)

	handler := withLogging(logger.Named("http"), router)

	logger.Info("Starting HTTP Server", zap.String("addr", httpAddr))
	if err := http.ListenAndServe(httpAddr, handler); err != nil {
		panic(err)
	}
}
