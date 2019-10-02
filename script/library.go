package script

import (
	"go.starlark.net/starlark"
)

// Library is a map of identifiers to underlying golang function implementations.
type Library map[string]Func

func (lib Library) Compile() starlark.StringDict {
	return lib.stringDict()
}

func (lib Library) stringDict() starlark.StringDict {
	symbols := starlark.StringDict{}
	if lib == nil {
		return symbols
	}

	for identifier, value := range lib {
		symbols[identifier] = value.builtin(identifier)
	}

	return symbols
}

// AddLibrary adds a library to the interpreter, exposing it to the execution environment.
func (i *Interpreter) AddLibrary(name string, lib Library) {
	i.libs[name] = lib.stringDict()
}
