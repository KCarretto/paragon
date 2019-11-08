package agent

import (
	"math/rand"
	"time"

	"go.uber.org/zap"
)

// A Transport wraps an implementation of Sender with metadata used for operation.
type Transport struct {
	Sender
	Log *zap.Logger

	Name     string
	Interval time.Duration
	Jitter   time.Duration
}

// Sleep for the transport's configured interval & jitter duration. If Jitter is non-zero, a random
// duration <= Jitter is added to the transports interval.
func (t Transport) Sleep() {
	// Random duration added to sleep delay [0, Jitter)
	var jitter time.Duration
	if max := t.Jitter.Nanoseconds(); max > 0 {
		jitter = time.Duration(rand.Int63n(max))
	}

	delay := t.Interval + jitter

	// Ratelimit log messages
	if delay > 1*time.Second {
		t.Log.Debug("Sleeping for transport interval + jitter", zap.Duration("delay", delay))
	}
	time.Sleep(delay)
}
