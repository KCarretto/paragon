package script

import (
	"io"
	"os"

	"go.starlark.net/starlark"
)

// Interpreter executes scripts within the environment of the preloaded libraries and builtins.
type Interpreter struct {
	builtins starlark.StringDict
	libs     map[string]starlark.StringDict
	output   io.Writer
}

// NewInterpreter initializes a new interpreter with sane defaults.
func NewInterpreter() *Interpreter {
	return &Interpreter{
		builtins: starlark.StringDict{},
		libs:     map[string]starlark.StringDict{},
		output:   os.Stdout,
	}
}

func (i *Interpreter) thread(name string, out io.Writer) *starlark.Thread {
	thread := &starlark.Thread{
		Name:  name,
		Print: i.printer(out),
		Load:  i.load,
	}

	return thread
}

func (i *Interpreter) printer(out io.Writer) func(t *starlark.Thread, msg string) {
	return func(_ *starlark.Thread, msg string) {
		// TODO: Handle error
		io.WriteString(out, msg)
	}
}

func (i *Interpreter) load(_ *starlark.Thread, module string) (starlark.StringDict, error) {
	lib, ok := i.libs[module]
	if !ok || lib == nil {
		return nil, ErrMissingLibrary
	}

	return lib, nil
}
