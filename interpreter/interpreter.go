package interpreter

import (
	"errors"
	"fmt"
	"io"

	"go.starlark.net/starlark"
)

// Interpreter TODO
type Interpreter struct {
	builtins starlark.StringDict
	stdlib   map[string]starlark.StringDict
}

// AddLibrary TODO
func (i Interpreter) AddLibrary(name string, src io.Reader) error {
	// Parse and validate syntax
	_, program, err := starlark.SourceProgram(name, src, i.builtins.Has)
	if err != nil {
		// TODO: Handle error
		return err
	}

	// Initialize & run top level declarations
	symbols, err := program.Init(i.thread(), i.builtins)
	if err != nil {
		// TODO: Handle error
		return err
	}

	// Add symbols to stdlib, making them available to load()
	i.stdlib[name] = symbols

	return nil
}

func (i Interpreter) thread() *starlark.Thread {
	thread := &starlark.Thread{
		// TODO: Thread IDs?
		Name:  "interpreter_thread",
		Print: i.print,
		Load:  i.load,
	}

	return thread
}

func (i Interpreter) print(_ *starlark.Thread, msg string) {
	fmt.Println(msg)
}

func (i Interpreter) load(_ *starlark.Thread, module string) (starlark.StringDict, error) {
	lib, ok := i.stdlib[module]
	if !ok || lib == nil {
		return nil, errors.New("module not found")
	}
	return lib, nil
}
