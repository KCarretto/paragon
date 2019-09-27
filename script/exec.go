package script

import (
	"context"

	"go.starlark.net/starlark"
	"go.uber.org/zap"
)

// EntryPoint defines the name of the method that will be called when running a script (if it exists)
// in addition to executing top level statements.
const EntryPoint = "main"

// Exec runs a script using the interpreter's execution environment.
func (i *Interpreter) Exec(ctx context.Context, logger *zap.Logger, script *Script) error {
	thread := i.thread(script.ID, logger)
	symbols, err := starlark.ExecFile(thread, script.ID, script, i.builtins)
	if err != nil {
		logger.Error("Failed to initialize script", zap.Error(err))
		return err
	}

	fn, ok := symbols[EntryPoint].(starlark.Callable)
	if !ok {
		return nil
	}

	res, err := starlark.Call(thread, fn, starlark.Tuple{}, []starlark.Tuple{})
	if err != nil {
		logger.Error("Failed to execute script", zap.Error(err))
		return err
	}

	if _, ok := res.(starlark.NoneType); !ok {
		logger.Info(res.String())
	}

	return nil
}
