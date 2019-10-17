// +build !dev,!debug

package main

import (
	"fmt"
	"io"

	"github.com/kcarretto/paragon/agent"
	"go.uber.org/zap"
)

func handleInterupt(logger *zap.Logger) int {
	logger.Info("Received interupt signal")
	return 0
}

func handleTerminate(logger *zap.Logger) int {
	logger.Error("Received terminate signal")
	return 0
}

func getLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println("logger initialization failed, logging disabled")
		logger = zap.NewNop()
	}
	return logger
}

func configureLogger(logger *zap.Logger, buf io.Writer) {}

func getTransports(logger *zap.Logger) (transports []agent.Transport) {
	// registry.Add()
	return
}
