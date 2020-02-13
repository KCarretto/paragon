package script

import (
	"fmt"
	"strings"

	"go.starlark.net/starlark"
)

// Library is a map of identifiers to underlying golang function implementations.
type Library map[string]Func

// String returns a description of the library, used to implement starlark.value
func (lib Library) String() string {
	return strings.Join(lib.AttrNames(), ", ")
}

// Freeze is a no-op since library methods are expected to be safe for concurrent use.
func (lib Library) Freeze() {}

// Type returns 'library' to indicate the type of the library within starlark.
func (lib Library) Type() string {
	return "library"
}

// Truth value of a library is True if it is non-empty
func (lib Library) Truth() starlark.Bool {
	if lib == nil {
		return starlark.False
	}
	return starlark.True
}

// Hash will error since the library type is not intended to be hashable.
func (lib Library) Hash() (uint32, error) {
	return 0, fmt.Errorf("library is unhashable")
}

// Attr enables dot expressions and returns the starlark method with the provided name if it exists.
func (lib Library) Attr(name string) (starlark.Value, error) {
	fn, ok := lib[name]
	if ok && fn != nil {
		return fn.builtin(name), nil
	}

	return nil, nil
}

// AttrNames returns the set of methods provided by the library.
func (lib Library) AttrNames() []string {
	keys := make([]string, 0, len(lib))
	for key := range lib {
		keys = append(keys, key)
	}
	return keys
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
