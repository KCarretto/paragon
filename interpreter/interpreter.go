package interpreter

import (
	"fmt"
	"io"

	"go.starlark.net/starlark"
)

// Interpreter TODO
type Interpreter struct {
	builtins starlark.StringDict
	libs     map[string]starlark.StringDict
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
		msgBytes := []byte(fmt.Sprintf("[%s] %s\n", t.Name, msg))
		output.Write(msgBytes)
	}
}

func (i *Interpreter) load(_ *starlark.Thread, module string) (starlark.StringDict, error) {
	lib, ok := i.libs[module]
	if !ok || lib == nil {
		return nil, ErrMissingLibrary
	}

	return lib, nil
}
