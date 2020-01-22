// +build nats

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"gocloud.dev/pubsub"
	_ "gocloud.dev/pubsub/gcppubsub"
	"github.com/kcarretto/paragon/pkg/event"
)

// NatsSubscriber implements Nats variant of the Subscriber interface.
type GCPSubscriber struct {
	wg     sync.WaitGroup
	ctx    context.Context
	cancel func()
}

func newSubscriber(ctx context.Context) (event.Subscriber, error) {
	newContext, cancel := context.WithCancel(ctx)
	return &GCPSubscriber{
		ctx:    newContext,
		cancel: cancel,
	}, nil
}

// Subscribes for events.
func (sub *GCPSubscriber) Subscribe(topic string, handler func(context.Context, []byte)) error {
	project := os.Getenv("GCP_PROJECT")
	if project == "" {
		return nil, fmt.Errorf("must set GCP_PROJECT environment variable to use GCP pubsub")
	}

	topic_uri := fmt.Sprintf("gcppubsub://projects/%s/topics/%s", project, topic)
	t, err := pubsub.OpenSubscription(ctx, topic_uri)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to the topic passed")
	}

	sub.wg.Add(1)

	go func() {
		defer sub.wg.Done()
		defer nc.Close()
		defer s.Shutdown(sub.ctx)

		for {
			msg, err := s.Receive(sub.ctx)
			if err != nil {
				log.Printf("[WARN] subscription channel has failed: %v", err)
				break
			}

			handler(sub.ctx, msg.Body)
			msg.Ack()
		}
	}()

	return nil
}

func (sub *GCPSubscriber) Close() error {
	sub.cancel()
	sub.wg.Wait()
	return nil
}
