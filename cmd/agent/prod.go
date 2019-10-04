// +build !dev,!local

package main

import (
	"fmt"
	"io"

	"github.com/kcarretto/paragon/transport"
	"github.com/kcarretto/paragon/transport/local"
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

func addTransports(logger *zap.Logger, receiver transport.PayloadWriter, registry *transport.Registry) {
	// TODO: Local HTTP
	registry.Add(transport.New(
		"local",
		local.New(receiver),
	))
}
