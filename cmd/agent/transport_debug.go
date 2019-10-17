// +build debug

package main

import (
	"time"

	"github.com/kcarretto/paragon/agent"
	"github.com/kcarretto/paragon/agent/debug"
	"go.uber.org/zap"
)

func transports(logger *zap.Logger) []agent.Transport {
	return []agent.Transport{
		agent.Transport{
			Name:     "Debug",
			Log:      logger.Named("debug"),
			Interval: 5 * time.Second,
			Jitter:   1 * time.Second,
			Sender: &debug.Sender{
				Log: logger.Named("debug"),
			},
		},
	}
}
