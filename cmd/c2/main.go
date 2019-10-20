package main

import (
	"context"
	"net/http"
	"os"

	"github.com/kcarretto/paragon/c2"

	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

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

	taskQueue, err := openSubscription(ctx, "tasks.queued")
	if err != nil {
		logger.Panic("Failed to subscribe to pubsub topic", zap.Error(err))
	}
	defer taskQueue.Shutdown(ctx)

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

	go listenForTasks(ctx, logger, srv, taskQueue)

	logger.Info("Starting HTTP Server", zap.String("addr", httpAddr))
	if err := http.ListenAndServe(httpAddr, srv); err != nil {
		panic(err)
	}
}
