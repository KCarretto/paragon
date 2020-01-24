package event

import (
	"context"
)

// A Broker publishes events and subscribes to them.
type Broker interface {
	Publisher
	Subscriber
}

// A Publisher broadcasts events to a topic.
type Publisher interface {
	Publish(context.Context, []byte) error
}

// A Subscriber registers event handlers for events received on a topic.
type Subscriber interface {
	Subscribe(topic string, handler func(context.Context, []byte)) error
}

// NewNopBroker returns a broker that does nothing.
func NewNopBroker() Broker {
	return nop(false)
}

type nop bool

func (nop) Publish(context.Context, []byte) error {
	return nil
}
func (nop) Subscribe(string, func(context.Context, []byte)) error {
	return nil
}
