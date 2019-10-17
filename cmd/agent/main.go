package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/kcarretto/paragon/agent"
	"github.com/pkg/profile"
	"go.uber.org/zap"
)

func runLoop() {
	pprof := profile.Start(profile.NoShutdownHook)

	// Initialize context
	ctx, cancel := context.WithCancel(context.Background())

	// Initialize logger
	logger := getLogger()

	// Initialize Agent
	paragon := &agent.Agent{
		Receiver: Receiver{
			ctx,
			logger.Named("agent.exec"),
		},
		Log:        logger.Named("agent"),
		Transports: getTransports(logger.Named("agent.transport")),
	}

	// Handle panic
	defer func() {
		if err := recover(); err != nil {
			logger.DPanic("Recovered from panic", zap.Reflect("error", err))
		}
	}()

	// Run agent
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		paragon.Run(ctx)
	}()

	// Listen for interupts
	sigint := make(chan os.Signal, 1)
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT)
	signal.Notify(sigterm, syscall.SIGTERM)

	// Wait for interupt
	select {
	case <-sigint:
		if exitCode := handleInterupt(logger); exitCode != 0 {
			pprof.Stop()
			os.Exit(exitCode)
		}
	case <-sigterm:
		if exitCode := handleTerminate(logger); exitCode != 0 {
			pprof.Stop()
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
