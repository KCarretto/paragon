package event

import (
	"context"
)

// A Publisher broadcasts events to a topic.
type Publisher interface {
	Publish(context.Context, []byte) error
}

// A Subscriber registers event handlers for events received on a topic.
type Subscriber interface {
	Subscribe(topic string, handler func(context.Context, []byte)) error
}
