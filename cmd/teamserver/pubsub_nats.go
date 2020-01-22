// +build nats

package main

import (
	"context"
	"fmt"
	"os"

	nats "github.com/nats-io/nats.go"
	"gocloud.dev/pubsub/natspubsub"
	"gocloud.dev/pubsub"
)

// NatsPublisher implements Nats variant of the Publisher interface.
type NatsPublisher struct {
	topic *pubsub.Topic
}

func newPublisher(ctx context.Context, topic string) (*NatsPublisher, error) {
	natsUrl := os.Getenv("NATS_URL")
	if natsUrl == "" {
		return nil, fmt.Errorf("must set NATS_URL environment variable to use Nats pubsub")
	}


	nc, err := nats.Connect(natsUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the nats server")
	}
	t, err := natspubsub.OpenTopic(nc, topic, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to the topic passed")
	}

	return &NatsPublisher{
		topic: t,
	}, nil
}

// Publish the event to stderr.
func (pub *NatsPublisher) Publish(ctx context.Context, data []byte) error {
	return pub.topic.Send(ctx, &pubsub.Message{
		Body: data,
	})
}
