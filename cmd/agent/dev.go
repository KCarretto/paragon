// +build dev

package main

import (
	"io"

	"github.com/kcarretto/paragon/transport"
	"github.com/kcarretto/paragon/transport/local"
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

func addTransports(receiver transport.PayloadWriter, registry *transport.Registry) {
	// TODO: Local HTTP
	registry.Add(transport.New(
		"local",
		local.New(receiver),
	))
}
