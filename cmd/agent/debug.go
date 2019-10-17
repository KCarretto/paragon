// +build debug

package main

import (
	"io"
	"time"

	"github.com/kcarretto/paragon/agent"
	"github.com/kcarretto/paragon/agent/debug"
	"go.uber.org/zap"
)

func handleInterupt(logger *zap.Logger) int {
	logger.Info("Received interupt signal")
	return 1
}

func handleTerminate(logger *zap.Logger) int {
	logger.Error("Received terminate signal")
	return 2
}

func getLogger() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return logger
}

func configureLogger(logger *zap.Logger, buf io.Writer) {}

func getTransports(logger *zap.Logger) []agent.Transport {
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
