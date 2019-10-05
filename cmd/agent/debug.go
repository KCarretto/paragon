// +build debug

package main

import (
	"io"

	"github.com/kcarretto/paragon/transport"
	"github.com/kcarretto/paragon/transport/debug"
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

func addTransports(logger *zap.Logger, receiver transport.PayloadWriter, registry *transport.Registry) {
	registry.Add(transport.New(
		"debug",
		&debug.Transport{
			Logger:        logger,
			PayloadWriter: receiver,
		},
	))
}
