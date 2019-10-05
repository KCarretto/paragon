package transport

import (
	"io"
	"time"
)

// An Option enables additional configuration of a transport.
type Option func(*Transport)

// Transport wraps an underlying io.Writer with metadata.
type Transport struct {
	io.Writer

	Name     string
	Priority int
	Interval time.Duration
	Jitter   time.Duration
}

// New creates and initializes a new Transport.
func New(name string, writer io.Writer, options ...Option) Transport {
	transport := Transport{
		Writer: writer,
		Name:   name,
	}

	for _, opt := range options {
		opt(&transport)
	}

	return transport
}

// SetPriority metadata for the transport.
func SetPriority(priority int) Option {
	return func(transport *Transport) {
		transport.Priority = priority
	}
}

// SetInterval metadata for the transport.
func SetInterval(interval time.Duration) Option {
	return func(transport *Transport) {
		transport.Interval = interval
	}
}

// SetJitter metadata for the transport.
func SetJitter(jitter time.Duration) Option {
	return func(transport *Transport) {
		transport.Jitter = jitter
	}
}
