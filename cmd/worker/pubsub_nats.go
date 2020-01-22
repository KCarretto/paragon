// +build nats

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	nats "github.com/nats-io/nats.go"
	"gocloud.dev/pubsub/natspubsub"
	"github.com/kcarretto/paragon/pkg/event"
)

// NatsSubscriber implements Nats variant of the Subscriber interface.
type NatsSubscriber struct {
	wg     sync.WaitGroup
	ctx    context.Context
	cancel func()
}

func newSubscriber(ctx context.Context) (event.Subscriber, error) {
	newContext, cancel := context.WithCancel(ctx)
	return &NatsSubscriber{
		ctx:    newContext,
		cancel: cancel,
	}, nil
}

// Subscribes for events.
func (sub *NatsSubscriber) Subscribe(topic string, handler func(context.Context, []byte)) error {
	natsUrl := os.Getenv("NATS_URL")
	if natsUrl == "" {
		return fmt.Errorf("must set NATS_URL environment variable to use Nats pubsub")
	}

	nc, err := nats.Connect(natsUrl)
	if err != nil {
		return fmt.Errorf("failed to connect to the nats server")
	}
	s, err := natspubsub.OpenSubscription(nc, topic, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe to the topic passed")
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

func (sub *NatsSubscriber) Close() error {
	sub.cancel()
	sub.wg.Wait()
	return nil
}
