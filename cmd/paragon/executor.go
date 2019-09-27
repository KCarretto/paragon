package main

import (
	"context"

	"github.com/kcarretto/paragon/agent"
	"github.com/kcarretto/paragon/script"
	"go.uber.org/zap"
)

// Executor is a wrappper around a starlark interpreter to enable it to run tasks as starlark scripts.
type Executor struct {
	*script.Interpreter
}

func (exec Executor) Exec(ctx context.Context, logger *zap.Logger, task agent.Task) {
	err := exec.Interpreter.Exec(ctx, logger, script.Script{
		Reader: task.Content,
		ID:     task.ID,
	})
	if err != nil {
		logger.Error("Task resulting in error", zap.Error(err))
	} else {
		logger.Info("Task completed successfully")
	}
}

func getExecutor() Executor {
	py := script.NewInterpreter()
	return Executor{py}
}
