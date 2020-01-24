// +build !gcp,!nats

package main

import (
	"context"

	"github.com/kcarretto/paragon/pkg/event"
)

func newPublisher(ctx context.Context, topic string) event.Publisher {
	return event.NewNopBroker()
}
