package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/kcarretto/paragon/pkg/agent"
	"github.com/kcarretto/paragon/pkg/agent/transport"
	"go.uber.org/zap"
)

// Run the agent, returns true when cancelled
func Run() bool {
	// Initialize context
	ctx, cancel := context.WithCancel(context.Background())

	// Initialize logger
	logger := newLogger().Named("agent")

	// Initialize Agent
	paragon := &agent.Agent{
		Log:          logger,
		MaxIdleTime:  3 * time.Second,
		TaskExecutor: Executor{},
		AgentMessageWriter: &transport.AgentMessageMultiWriter{
			Transports: transports(logger.Named("transport")),
		},
	}
	paragon.OnRun = func(){
		paragon.TaskExecutor = Executor{paragon.Metadata}
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
		for {
			if err := ctx.Err(); err != nil {
				return
			}

			paragon.Run(ctx)
		}
	}()

	// Listen for interupts
	sigint := make(chan os.Signal, 1)
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT)
	signal.Notify(sigterm, syscall.SIGTERM)

	// Wait for interupt
	select {
	case <-sigint:
		logger.Info("Received interupt signal")
	case <-sigterm:
		logger.Error("Received terminate signal")
	}

	// Wait for threads to finish
	cancel()
	wg.Wait()

	return true
}

func main() {
	// Enable profiling for the run if compiled
	pprof := startProfile()
	defer pprof.Stop()

	interupted := false
	for !interupted {
		interupted = Run()
	}
}
