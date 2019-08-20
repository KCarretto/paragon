package interpreter

import (
	"io"

	"go.starlark.net/starlark"
)

// Retval TODO
type Retval interface{}

// Func provides a simple wrapper that enables a golang function to be exposed to starlark.
type Func func(args ArgParser, output io.Writer) (Retval, error)

// toBuiltin wraps the Func and returns a starlark builtin
func (fn Func) toBuiltin(name string, output io.Writer) *starlark.Builtin {
	return starlark.NewBuiltin(
		name,
		func(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			parser, err := getParser(args, kwargs)
			if err != nil {
				return nil, err
			}

			retval, err := fn(parser, output)
			if err != nil {
				return nil, err
			}

			val, err := convertToStarlark(retval)
			if err != nil {
				return nil, err
			}

			return val, nil
		},
	)
}
