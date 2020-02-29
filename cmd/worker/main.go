package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/kcarretto/paragon/graphql"
	"github.com/kcarretto/paragon/graphql/models"
	"github.com/kcarretto/paragon/pkg/cdn"
	"github.com/kcarretto/paragon/pkg/event"
	"github.com/kcarretto/paragon/pkg/worker"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	logger = logger.Named("svc.worker")

	logger.Debug("Initializing worker")

	teamserverURL := "http://127.0.0.1:80"
	if url := os.Getenv("TEAMSERVER_URL"); url != "" {
		teamserverURL = url
	}
	graph := graphql.Client{
		URL:     fmt.Sprintf("%s/%s", teamserverURL, "graphql"),
		Service: "pg-worker",
	}
	// TODO: Fix this, but this initial call is to register the service
	graph.ClaimTasks(context.Background(), models.ClaimTasksRequest{
		PrimaryIP: nil,
	})

	cdnURL := teamserverURL
	if url := os.Getenv("CDN_URL"); url != "" {
		cdnURL = url
	}
	cdn := cdn.Client{URL: cdnURL}

	w := &worker.Worker{
		Uploader:   cdn,
		Downloader: cdn,
		Graph:      graph,
	}

	sub, _ := newSubscriber(context.Background())
	if closer, ok := sub.(io.Closer); ok {
		defer closer.Close()
	}

	topic := os.Getenv("PUB_TOPIC")
	if topic == "" {
		logger.Warn("[WARN] No PUB_TOPIC environment variable set to publish events, is this a mistake?")
	}

	var wg sync.WaitGroup
	defer wg.Wait()

	taskHandler := func(ctx context.Context, data []byte) {
		defer func() {
			if fatalErr := recover(); fatalErr != nil {
				logger.Error("Recovered from panic!", zap.Reflect("err", fatalErr))
			}
		}()

		logger.Info("Message Received", zap.String("data", string(data)))

		var taskQueuedEvent event.TaskQueued
		if err := json.Unmarshal(data, &taskQueuedEvent); err != nil {
			logger.Error("Failed to parse event json", zap.Error(err))
			return
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			w.HandleTaskQueued(ctx, taskQueuedEvent)
		}()
	}

	if err := sub.Subscribe(topic, taskHandler); err != nil {
		panic(err)
	}

	logger.Info("Worker Initialized",
		zap.String("pub_topic", topic),
		zap.String("teamserver_url", teamserverURL),
		zap.String("cdn_url", cdnURL),
	)

	// Listen for interupts
	sigint := make(chan os.Signal, 1)
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT)
	signal.Notify(sigterm, syscall.SIGTERM)

	// Wait for interupt
	select {
	case <-sigint:
		logger.Info("Received interupt signal")
	case <-sigterm:
		logger.Error("Received terminate signal")
	}
}
