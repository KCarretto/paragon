package script

import (
	"io"
	"os"

	"go.starlark.net/starlark"
	"go.uber.org/zap"
)

// Interpreter executes scripts within the environment of the preloaded libraries and builtins.
type Interpreter struct {
	builtins starlark.StringDict
	libs     map[string]starlark.StringDict
	output   io.Writer
	logger   *zap.Logger
}

// NewInterpreter initializes a new interpreter with sane defaults.
func NewInterpreter() *Interpreter {
	return &Interpreter{
		builtins: starlark.StringDict{},
		libs:     map[string]starlark.StringDict{},
		output:   os.Stdout,
		logger:   zap.NewNop(),
	}
}

func (i *Interpreter) thread(name string, logger *zap.Logger) *starlark.Thread {
	thread := &starlark.Thread{
		Name:  name,
		Print: i.printer(logger.With(zap.String("thread_name", name))),
		Load:  i.load,
	}

	return thread
}

func (i *Interpreter) printer(logger *zap.Logger) func(t *starlark.Thread, msg string) {
	return func(_ *starlark.Thread, msg string) {
		i.logger.Info(msg)
	}
}

func (i *Interpreter) load(_ *starlark.Thread, module string) (starlark.StringDict, error) {
	lib, ok := i.libs[module]
	if !ok || lib == nil {
		return nil, ErrMissingLibrary
	}

	return lib, nil
}
