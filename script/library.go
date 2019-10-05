package script

import (
	"go.starlark.net/starlark"
)

// Library is a map of identifiers to underlying golang function implementations.
type Library map[string]Func

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
