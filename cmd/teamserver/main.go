package main

import (
	"context"
	"log"
	"os"

	"github.com/kcarretto/paragon/pkg/teamserver"
)

const DEFAULT_TOPIC = "TOPIC"

func main() {
	topic := os.Getenv("PUB_TOPIC")
	if topic == "" {
		log.Println("You did not set PUB_TOPIC environment variable to publish events, is this a mistake?")
	}

	ctx := context.Background()

	client := getClient(ctx)
	defer client.Close()

	pub, err := newPublisher(ctx, topic)
	if err != nil {
		panic("You failed to create a publisher for the assigned topic... exiting")
	}
	defer pub.topic.Shutdown(ctx)

	server := teamserver.Server{
		Log:       newLogger(),
		EntClient: client,
		Publisher: pub,
	}
	log.Println("Starting teamserver")
	server.Run()
}
