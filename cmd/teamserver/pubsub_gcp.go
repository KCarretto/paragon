// +build gcp

package main

import (
	"context"
	"fmt"
	"os"

	"gocloud.dev/pubsub"
	_ "gocloud.dev/pubsub/gcppubsub"
)

// NatsPublisher implements GCP variant of the Publisher interface.
type GCPPublisher struct {
	topic *pubsub.Topic
}

func newPublisher(ctx context.Context, topic string) (*GCPPublisher, error) {
	project := os.Getenv("GCP_PROJECT")
	if project == "" {
		return nil, fmt.Errorf("must set GCP_PROJECT environment variable to use GCP pubsub")
	}

	topic_uri := fmt.Sprintf("gcppubsub://projects/%s/topics/%s", project, topic)
	t, err := pubsub.OpenTopic(ctx, topic_uri)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to the topic passed")
	}

	return &GCPPublisher{
		topic: t,
	}, nil
}

func (pub *GCPPublisher) Publish(ctx context.Context, data []byte) error {
	return pub.topic.Send(ctx, &pubsub.Message{
		Body: data,
	})
}
