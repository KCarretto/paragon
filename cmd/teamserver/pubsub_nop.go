// +build !gcp

package main

import (
	"context"
	"log"

	"gocloud.dev/pubsub"
	_ "gocloud.dev/pubsub/mempubsub"
)

// NOPPublisher implements a no-op variant of the Publisher interface. It publishes events to stderr.
type NOPPublisher struct {
	topic *pubsub.Topic
}

func newPublisher(ctx context.Context, topic string) (*NOPPublisher, error) {
	return &NOPPublisher{
		topic: nil,
	}, nil
}

// Publish the event to stderr.
func (pub *NOPPublisher) Publish(ctx context.Context, data []byte) error {
	log.Println(string(data))
	return nil
}
