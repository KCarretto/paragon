package transport

import (
	"io"
	"time"

	"go.uber.org/zap"
)

// A Factory is used to initialize a new transport.
type Factory func(*zap.Logger, Tasker) (io.WriteCloser, error)

// A Tasker is provided to a Transport to enable it to send incoming tasks to the agent.
type Tasker interface {
	QueueTask(id string, content io.Reader)
}

// An Option enables additional configuration of a transport Meta.
type Option func(*Meta)

// Meta holds metadata about a transport and a factory method to initialize it.
type Meta struct {
	Name     string
	Priority int
	Interval time.Duration
	Jitter   time.Duration
	Factory  Factory
}

// SetPriority metadata for the transport.
func SetPriority(priority int) Option {
	return func(meta *Meta) {
		meta.Priority = priority
	}
}

// SetInterval metadata for the transport.
func SetInterval(interval time.Duration) Option {
	return func(meta *Meta) {
		meta.Interval = interval
	}
}

// SetJitter metadata for the transport.
func SetJitter(jitter time.Duration) Option {
	return func(meta *Meta) {
		meta.Jitter = jitter
	}
}
