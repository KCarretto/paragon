package main

import (
	"context"

	"github.com/kcarretto/paragon/teamserver"
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

	claimedSubscription, err := openSubscription(ctx, "tasks.claimed")
	if err != nil {
		panic(err)
	}
	defer claimedSubscription.Shutdown(ctx)

	executedSubscription, err := openSubscription(ctx, "tasks.executed")
	if err != nil {
		panic(err)
	}
	defer executedSubscription.Shutdown(ctx)

	server := teamserver.Server{
		Log:                  newLogger(),
		EntClient:            client,
		QueuedTopic:          queuedTopic,
		ClaimedSubscription:  claimedSubscription,
		ExecutedSubscription: executedSubscription,
	}
	server.Run()
}
