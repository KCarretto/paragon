package interpreter

import (
	"go.starlark.net/starlark"
)

func convertToStarlark(value interface{}) (starlark.Value, error) {
	// TODO
	return nil, nil
}

// Retval TODO
type Retval interface{}

// Func provides a simple wrapper that enables a golang function to be exposed to starlark.
type Func func(args ArgParser) (Retval, error)

// toBuiltin wraps the Func and returns a starlark builtin
func (fn Func) toBuiltin(name string) *starlark.Builtin {
	return starlark.NewBuiltin(
		name,
		func(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			parser, err := getParser(args, kwargs)
			if err != nil {
				return nil, err
			}

			retval, err := fn(parser)
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
