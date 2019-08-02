package main

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/kcarretto/paragon/internal/agent"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Debug("Exiting agent")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	agent := agent.New(
		agent.WithLogger(logger),
	)
	agent.Run(ctx)
}
