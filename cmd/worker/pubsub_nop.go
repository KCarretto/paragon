// +build !gcp,!nats

package main

import (
	"context"

	"github.com/kcarretto/paragon/pkg/event"
)

// NOPSubscriber implements Nats variant of the Subscriber interface.
type NOPSubscriber struct {
}

func newSubscriber(_ context.Context) (event.Subscriber, error) {
	return &NOPSubscriber{}, nil
}

// Subscribe for events.
func (sub *NOPSubscriber) Subscribe(topic string, handler func(context.Context, []byte)) error {
	return nil
}
