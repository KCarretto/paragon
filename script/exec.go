package script

import (
	"context"
	"io"

	"go.starlark.net/starlark"
)

// EntryPoint defines the name of the method that will be called when running a script (if it exists)
// in addition to executing top level statements.
const EntryPoint = "main"

// Exec runs a script using the interpreter's execution environment.
func (i *Interpreter) Exec(ctx context.Context, script Script, output io.Writer) error {
	thread := i.thread(script.ID, output)
	symbols, err := starlark.ExecFile(thread, script.ID, script, i.builtins)
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
		// TODO: Handle error
		io.WriteString(output, res.String())
	}

	return nil
}
