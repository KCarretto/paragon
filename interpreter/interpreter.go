package interpreter

import (
	"fmt"
	"io"
	"os"

	"go.starlark.net/starlark"
)

// Interpreter TODO
type Interpreter struct {
	builtins starlark.StringDict
	libs     map[string]starlark.StringDict
	output   io.Writer
}

// New TODO
func New() *Interpreter {
	return &Interpreter{
		builtins: starlark.StringDict{},
		libs:     map[string]starlark.StringDict{},
		output:   os.Stdout,
	}
}

func (i *Interpreter) thread(name string, output io.Writer) *starlark.Thread {
	thread := &starlark.Thread{
		Name:  name,
		Print: i.printer(output),
		Load:  i.load,
	}

	return thread
}

func (i *Interpreter) printer(output io.Writer) func(t *starlark.Thread, msg string) {
	return func(t *starlark.Thread, msg string) {
		// TODO: Format with logger + timestamp?
		io.WriteString(output, fmt.Sprintf("[%s] %s\n", t.Name, msg))
	}
}

func (i *Interpreter) load(_ *starlark.Thread, module string) (starlark.StringDict, error) {
	lib, ok := i.libs[module]
	if !ok || lib == nil {
		return nil, ErrMissingLibrary
	}

	return lib, nil
}
