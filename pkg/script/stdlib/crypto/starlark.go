package crypto

import (
	"fmt"

	"github.com/kcarretto/paragon/pkg/script"
	"go.starlark.net/starlark"
)

// Key is a basic type alias for string that will let us provide linting within scripts
// WARNING: This is NOT how you should normally handle keys in go.
type Key struct {
	value string
}

/*
 * Starlark.Value methods
 */

// String returns the name of the file.
func (k Key) String() string {
	return k.value
}

// Freeze is a no-op since the underlying file is safe for concurrent use.
func (k Key) Freeze() {}

// Type returns 'file' to indicate the type of the file within starlark.
func (k Key) Type() string {
	return "key"
}

// Truth value of a file is True if the file is non-nil
func (k Key) Truth() starlark.Bool {
	if k.value == "" {
		return starlark.False
	}
	return starlark.True
}

// Hash will error since the file type is not intended to be hashable.
func (k Key) Hash() (uint32, error) {
	return 0, fmt.Errorf("key type is unhashable")
}

// ParseParam from starlark input
func ParseParam(parser script.ArgParser, index int) (Key, error) {
	val, err := parser.GetParam(index)
	if err != nil {
		return Key{""}, err
	}

	k, ok := val.(Key)
	if !ok {
		return Key{""}, fmt.Errorf("%w: expected key type", script.ErrInvalidArgType)
	}

	return k, nil
}
