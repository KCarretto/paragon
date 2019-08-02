package agent

import (
	"math/rand"
	"time"

	"go.uber.org/zap"
)

func (a *agent) sleep() {
	if a.interval <= 0 {
		a.logger.Debug(
			"Agent interval <= 0, skipping sleep period",
			zap.Duration("interval", a.interval),
		)
		return
	}

	// Seed RNG
	rand.Seed(time.Now().UTC().UnixNano())

	var duration time.Duration
	{
		maxJitter := a.jitter.Nanoseconds()
		if maxJitter > 0 {
			jitter := time.Duration(rand.Int63n(maxJitter) - rand.Int63n(maxJitter))
			duration = a.interval + jitter
		} else {
			duration = a.interval
		}
	}

	// Sleep for calculated duration
	a.logger.Debug("Agent sleeping", zap.Duration("duration", duration))
	time.Sleep(duration)
}
