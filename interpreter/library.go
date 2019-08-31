package interpreter

import (
	"go.starlark.net/starlark"
)

// Library is a map of identifiers to underlying golang function implementations.
type Library map[string]Func

// AddLibrary exposes golang bindings to the starlark interpreter.
func (i *Interpreter) AddLibrary(name string, bindings Library) error {
	if i.libs == nil {
		i.libs = map[string]starlark.StringDict{}
	}

	lib := starlark.StringDict{}
	for binding, fn := range bindings {
		lib[binding] = fn.toBuiltin(binding, i.output)
	}

	i.libs[name] = lib

	return nil
}
