// +build nats

package main

import (
	"context"
	"fmt"
	"os"
	// "time"
	nats "github.com/nats-io/nats.go"
	"gocloud.dev/pubsub"
	"gocloud.dev/pubsub/natspubsub"
)

// NatsPublisher implements Nats variant of the Publisher interface.
type NatsPublisher struct {
	topic *pubsub.Topic
}

func newPublisher(ctx context.Context, topic string) *NatsPublisher {
	natsUrl := os.Getenv("NATS_URL")
	if natsUrl == "" {
		panic(fmt.Errorf("must set NATS_URL environment variable to use Nats pubsub"))
	}

	nc, err := nats.Connect(natsUrl)
	if err != nil {
		// time.Sleep(10 * time.Minute)
		panic(fmt.Errorf("failed to connect to the nats server %s", natsUrl))
	}
	t, err := natspubsub.OpenTopic(nc, topic, nil)
	if err != nil {
		panic(fmt.Errorf("failed to subscribe to the topic passed"))
	}

	return &NatsPublisher{
		topic: t,
	}
}

// Publish the event to stderr.
func (pub *NatsPublisher) Publish(ctx context.Context, data []byte) error {
	return pub.topic.Send(ctx, &pubsub.Message{
		Body: data,
	})
}
