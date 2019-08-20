package interpreter

import (
	"io"

	"go.starlark.net/starlark"
)

// EntryPoint defines the name of the method that will be called when running a script (if it exists)
// in addition to executing top level statements.
const EntryPoint = "main"

// Execute a script using the interpreter.
func (i *Interpreter) Execute(script *Script, output io.Writer) error {
	thread := i.thread(script.Name, output)
	symbols, err := starlark.ExecFile(thread, script.Name, script, i.builtins)
	if err != nil {
		return err
	}

	fn, ok := symbols[EntryPoint].(starlark.Callable)
	if !ok {
		return nil
	}

	res, err := starlark.Call(thread, fn, starlark.Tuple{}, []starlark.Tuple{})
	if err != nil {
		return err
	}
	output.Write([]byte(res.String()))

	return nil
}
