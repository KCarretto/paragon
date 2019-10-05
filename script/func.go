package script

import (
	"go.starlark.net/starlark"
)

// Retval TODO
type Retval interface{}

// Func provides a simple wrapper that enables a golang function to be exposed to starlark.
type Func func(args ArgParser) (Retval, error)

// builtin wraps the Func and returns a starlark builtin
func (fn Func) builtin(name string) *starlark.Builtin {
	return starlark.NewBuiltin(
		name,
		func(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
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
