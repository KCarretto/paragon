package main

import (
	"os"
	"os/signal"
	"syscall"
)

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
