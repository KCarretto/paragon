// +build !gcp

package main

// import (
// 	"context"
// 	"fmt"

// 	"gocloud.dev/pubsub"
// 	_ "gocloud.dev/pubsub/mempubsub"
// )

// func openTopic(ctx context.Context, topic string) (*pubsub.Topic, error) {
// 	uri := fmt.Sprintf("mem://%s", topic)
// 	return pubsub.OpenTopic(ctx, uri)
// }

// func openSubscription(ctx context.Context, topic string) (*pubsub.Subscription, error) {
// 	uri := fmt.Sprintf("mem://%s", topic)
// 	return pubsub.OpenSubscription(ctx, uri)
// }
