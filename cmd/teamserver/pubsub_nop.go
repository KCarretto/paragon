// +build !gcp

package main

import (
	"context"
	"log"

	"gocloud.dev/pubsub"
	_ "gocloud.dev/pubsub/mempubsub"
)

type NOPPublisher struct {
	topic *pubsub.Topic
}

func newPublisher(ctx context.Context, topic string) (*NOPPublisher, error) {
	return &NOPPublisher{
		topic: nil,
	}, nil
}

func (pub *NOPPublisher) Publish(ctx context.Context, data []byte) error {
	log.Println(string(data))
	return nil
}
