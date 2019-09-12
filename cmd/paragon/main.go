package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/kcarretto/paragon/agent"
	"go.uber.org/zap"
)

type Runner string

func (runner Runner) Run(context.Context, *zap.Logger, agent.Task) {}

func main() {
	interupts := make(chan os.Signal, 1)
	signal.Notify(interupts, syscall.SIGINT, syscall.SIGTERM)

	logger := getLogger()

	go Run(logger)

	switch <-interupts {
	case syscall.SIGINT:
		logger.Info("Received interupt signal")
	case syscall.SIGTERM:
		logger.Error("Received terminate signal")
	}

}
