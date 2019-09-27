package main

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	interupts := make(chan os.Signal, 1)
	signal.Notify(interupts, syscall.SIGINT, syscall.SIGTERM)

	logger := getLogger()

	bot := NewAgent()
	addTransports(bot)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.DPanic("Recovered from panic", zap.Reflect("error", err))
			}
		}()

		err := bot.Run()
		defer func() {
			if err := bot.Close(); err != nil {
				panic(err)
			}
		}()

		if err != nil {
			// TODO: Handle ErrNoTransport
			logger.DPanic("Run encountered an error", zap.Error(err))
		}
	}()

	switch <-interupts {
	case syscall.SIGINT:
		logger.Info("Received interupt signal")
	case syscall.SIGTERM:
		logger.Error("Received terminate signal")
	}

}
