package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func runLoop() {
	// Initialize context
	ctx, cancel := context.WithCancel(context.Background())

	// Initialize logger
	logger := getLogger()

	// Handle panic
	defer func() {
		if err := recover(); err != nil {
			logger.DPanic("Recovered from panic", zap.Reflect("error", err))
		}
	}()

	// Run agent
	var wg sync.WaitGroup
	wg.Add(1)
	go func(logger *zap.Logger) {
		defer wg.Done()
		run(ctx, logger)
	}(logger.With(zap.Time("agent_loop_start", time.Now())))

	// Listen for interupts
	sigint := make(chan os.Signal, 1)
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT)
	signal.Notify(sigterm, syscall.SIGTERM)

	// Wait for interupt
	select {
	case <-sigint:
		if exitCode := handleInterupt(logger); exitCode != 0 {
			os.Exit(exitCode)
		}
	case <-sigterm:
		if exitCode := handleTerminate(logger); exitCode != 0 {
			os.Exit(exitCode)
		}
	}

	// After interupt, wait for threads to finish
	cancel()
	wg.Wait()

}

func main() {
	for {
		runLoop()
	}
}
