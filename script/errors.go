package script

import (
	"errors"
)

var (
	// ErrMalformattedKwarg occurs when a keyword argument passed to a starlark function is malformatted.
	ErrMalformattedKwarg = errors.New("malformatted keyword argument for method call")

	// ErrMissingKwarg occurs when a keyword argument is missing from a starlark function call.
	ErrMissingKwarg = errors.New("missing keyword argument for method call")

	// ErrMissingArg occurs when an argument is missing from a starlark function call.
	ErrMissingArg = errors.New("missing argument for method call")

	// ErrInvalidArgType occurs when an argument provided to a starlark function call has the wrong type.
	ErrInvalidArgType = errors.New("invalid argument type provided to method")

	// ErrMissingLibrary occurs when attempting to load a library that doesn't exist.
	ErrMissingLibrary = errors.New("could not find library to load")

	// ErrInvalidTypeConversion occurs when converting a golang type to a starlark.Value fails.
	ErrInvalidTypeConversion = errors.New("could not convert golang value to starlark.Value")
)
