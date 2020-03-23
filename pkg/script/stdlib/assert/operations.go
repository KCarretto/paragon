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

// Equal will check if two values are equal. This function will result in a fatal error if the assertion is incorrect.
//
//go:generate go run ../gendoc.go -lib assert -func equal -param expected@starlark.Value -param actual@starlark.Value -doc "Equal will check if two values are equal. This function will result in a fatal error if the assertion is incorrect."
//
// @callable:	assert.equal
// @param:		expected @starlark.Value
// @param:		actual	 @starlark.Value
//
// @usage:		assert.equal("expected string", some_string)
func Equal(expected starlark.Value, actual starlark.Value) error {
	if expected != actual {
		return fmt.Errorf("assertion failed: values are not equal. Expected: %v, Got: %v", expected, actual)
	}

	return nil
}

func equal(parser script.ArgParser) (script.Retval, error) {
	expectedVal, err := parser.GetParam(0)
	if err != nil {
		return nil, err
	}
	actualVal, err := parser.GetParam(1)
	if err != nil {
		return nil, err
	}

	retErr := Equal(expectedVal, actualVal)
	return nil, retErr
}
