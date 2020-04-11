package script

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"go.starlark.net/starlark"
)

// EntryPoint defines the name of the method that will be called when running a script (if it exists)
// in addition to executing top level statements.
const EntryPoint = "main"

// A Script provides metadata and instructions to be executed by the interpreter.
type Script struct {
	io.Reader // Read() instructions to execute
	io.Writer // Write() execution output
	ID        string
	Builtins  map[string]Func
	Libraries map[string]Library

	thread  *starlark.Thread
	symbols starlark.StringDict
}

// Exec calls the entrypoint method of a script
func (script *Script) Exec(ctx context.Context) error {
	_, err := script.Call(EntryPoint, starlark.Tuple{})
	return err
}

// Call a function defined by the script.
// This will execute all top level statements before calling the function.
func (script *Script) Call(funcName string, args starlark.Tuple, kwargs ...starlark.Tuple) (starlark.Value, error) {
	thread := script.newThread()
	symbols, err := script.parseSymbols(thread)
	if err != nil {
		return starlark.None, err
	}

	fn, ok := symbols[funcName].(starlark.Callable)
	if !ok {
		return starlark.None, nil
	}

	res, err := starlark.Call(thread, fn, args, kwargs)
	if err != nil {
		if evalErr, ok := err.(*starlark.EvalError); ok {
			return starlark.None, fmt.Errorf("%s", evalErr.Backtrace())
		}
		return starlark.None, err
	}

	if _, ok := res.(starlark.NoneType); !ok {
		thread.Print(thread, res.String())
	}

	return res, nil
}

func (script *Script) parseSymbols(thread *starlark.Thread) (starlark.StringDict, error) {
	if script.symbols == nil {
		builtins := script.compilePredeclared()
		symbols, err := starlark.ExecFile(thread, script.ID, script.Reader, builtins)
		if err != nil {
			if evalErr, ok := err.(*starlark.EvalError); ok {
				return nil, fmt.Errorf("%s", evalErr.Backtrace())
			}
			return nil, err
		}
		script.symbols = symbols
	}

	return script.symbols, nil
}

func (script *Script) compilePredeclared() starlark.StringDict {
	builtins := make(starlark.StringDict)

	for name, fn := range script.Builtins {
		builtins[name] = fn.builtin(name)
	}
	for name, lib := range script.Libraries {
		builtins[name] = lib
	}

	return builtins
}

func (script *Script) newThread() *starlark.Thread {
	if script.thread == nil {
		name := "Renegade"
		if script.ID != "" {
			name = script.ID
		}
		script.thread = &starlark.Thread{
			Name: name,
			Print: func(_ *starlark.Thread, msg string) {
				// TODO: Handle error
				fmt.Fprintf(script.Writer, "%s\n", msg)
			},
			Load: func(_ *starlark.Thread, module string) (starlark.StringDict, error) {
				lib, ok := script.Libraries[module]
				if !ok || lib == nil {
					// TODO: Wrap error with module name & script id.
					return nil, ErrMissingLibrary
				}

				return lib.stringDict(), nil
			},
		}
	}

	return script.thread
}

// New initializes and returns a script with the provided contents.
func New(id string, content io.Reader, options ...Option) *Script {
	script := &Script{
		Reader:    content,
		Writer:    ioutil.Discard,
		ID:        id,
		Builtins:  map[string]Func{},
		Libraries: map[string]Library{},
	}

	for _, opt := range options {
		opt(script)
	}

	return script
}
