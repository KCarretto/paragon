package main

import (
	"context"

	"github.com/kcarretto/paragon/pkg/teamserver"
)

func main() {
	ctx := context.Background()

	client := getClient(ctx)
	defer client.Close()

	queuedTopic, err := openTopic(ctx, "tasks.queued")
	if err != nil {
		panic(err)
	}
	defer queuedTopic.Shutdown(ctx)

	server := teamserver.Server{
		Log:         newLogger(),
		EntClient:   client,
		QueuedTopic: queuedTopic,
	}
	server.Run()
}
