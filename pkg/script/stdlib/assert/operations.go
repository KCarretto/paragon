package assert

import (
	"fmt"

	"github.com/kcarretto/paragon/pkg/script"
	"go.starlark.net/starlark"
)

// NoError will check if the passed value is a starlark.NoneType, if not it will error out the
// script. This function may cause a fatal error if the assertion is incorrect.
//
//go:generate go run ../gendoc.go -lib assert -func noError -param err@starlark.Value -doc "NoError will check if the passed value is a starlark.NoneType, if not it will error out the script.  This function may cause a fatal error if the assertion is incorrect."
//
// @callable:	assert.noError
// @param:		err	@starlark.Value
//
// @usage:		assert.noError(err)
func NoError(err starlark.Value) error {
	_, ok := err.(starlark.NoneType)
	if !ok {
		errString, ok := err.(starlark.String)
		if !ok {
			return fmt.Errorf("assertion failed: error is not None")
		}
		return fmt.Errorf("assertion failed: error is not None: %v", errString)
	}
	return nil
}

func noError(parser script.ArgParser) (script.Retval, error) {
	errVal, err := parser.GetParam(0)
	if err != nil {
		return nil, err
	}

	retErr := NoError(errVal)
	return nil, retErr
}
