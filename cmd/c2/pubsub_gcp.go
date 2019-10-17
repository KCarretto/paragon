// +build gcp

package main

// import (
// 	"context"
// 	"fmt"
// 	"os"

// 	"gocloud.dev/pubsub"
// 	_ "gocloud.dev/pubsub/gcppubsub"
// )

// func getURI(topic string) (string, error) {
// 	project := os.Getenv("GCP_PROJECT")
// 	if project == "" {
// 		return "", fmt.Errorf("must set GCP_PROJECT environment variable to use GCP pubsub")
// 	}

// 	uri := fmt.Sprintf("gcppubsub://projects/%s/topics/%s", project, topic)
// 	return uri, nil
// }

// func openTopic(ctx context.Context, topic string) (*pubsub.Topic, error) {
// 	uri, err := getURI(topic)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return pubsub.OpenTopic(ctx, uri)
// }

// func openSubscription(ctx context.Context, topic string) (*pubsub.Subscription, error) {
// 	uri, err := getURI(topic)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return pubsub.OpenSubscription(ctx, uri)
// }
