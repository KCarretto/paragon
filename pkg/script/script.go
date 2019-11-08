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
}

// Exec parses input from the script reader and executes it. It will also invoke the EntryPoint
// method if one is available.
func (script Script) Exec(ctx context.Context) error {
	builtins := script.compilePredeclared()
	thread := script.newThread()

	symbols, err := starlark.ExecFile(thread, script.ID, script.Reader, builtins)
	if err != nil {
		// TODO: Better error type
		return err
	}

	fn, ok := symbols[EntryPoint].(starlark.Callable)
	if !ok {
		return nil
	}

	res, err := starlark.Call(thread, fn, starlark.Tuple{}, []starlark.Tuple{})
	if err != nil {
		// TODO: Better error type
		return err
	}

	if _, ok := res.(starlark.NoneType); !ok {
		thread.Print(thread, res.String())
	}

	return nil
}

func (script Script) compilePredeclared() (builtins starlark.StringDict) {
	for name, fn := range script.Builtins {
		builtins[name] = fn.builtin(name)
	}
	return
}

func (script Script) newThread() *starlark.Thread {
	return &starlark.Thread{
		Name: script.ID,
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

// New initializes and returns a script with the provided contents.
func New(id string, content io.Reader, options ...Option) Script {
	script := Script{
		content,
		ioutil.Discard,
		id,
		map[string]Func{},
		map[string]Library{},
	}

	for _, opt := range options {
		opt(&script)
	}

	return script
}
