//+build !development

package main

import (
	"github.com/kcarretto/paragon/agent"
	"go.uber.org/zap"
)

func getLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	return logger
}

func runLoop(logger *zap.Logger) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Recovered from panic", zap.Reflect("error", err))
		}
	}()

	agent := agent.New(
		logger,
		Executor("Stuff"),
	)
	defer func() {
		if err := agent.Close(); err != nil {
			logger.Error("Encountered error closing agent", zap.Error(err))
		}
	}()

	agent.Run()
}

func Run(logger *zap.Logger) {
	for {
		runLoop(logger)
	}
}
