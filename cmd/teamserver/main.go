package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/kcarretto/paragon/pkg/teamserver"
)

// DefaultTopic is the default pubsub topic where events will be published. It can be configured
// by setting the PUB_TOPIC environment variable.
const DefaultTopic = "TOPIC"

func main() {
	ctx := context.Background()

	topic := os.Getenv("PUB_TOPIC")
	if topic == "" {
		log.Println("[WARN] No PUB_TOPIC environment variable set to publish events, is this a mistake?")
	}

	client := getClient(ctx)
	defer client.Close()

	pub, err := newPublisher(ctx, topic)
	if err != nil {
		panic(fmt.Errorf("You failed to create a publisher for the assigned topic: %+v", err))
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
