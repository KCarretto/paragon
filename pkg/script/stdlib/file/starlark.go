package file

import (
	"fmt"

	"github.com/kcarretto/paragon/pkg/script"

	"github.com/spf13/afero"
	"go.starlark.net/starlark"
)

// File descriptor that can be used with the file library and implements starlark.Value.
type File struct {
	afero.Fs

	Path string
}

// Name returns the path used to create the file descriptor.
func (f *File) Name() string {
	return f.Path
}

// String returns the name of the file.
func (f *File) String() string {
	return f.Path
}

// Type returns 'file' to indicate the type of the file within starlark.
func (f *File) Type() string {
	return "file"
}

// Freeze is a no-op since the underlying file is safe for concurrent use.
func (f *File) Freeze() {}

// Hash will error since the file type is not intended to be hashable.
func (f *File) Hash() (uint32, error) {
	return 0, fmt.Errorf("file type is unhashable")
}

// Truth value of a file is True if the file is non-nil
func (f *File) Truth() starlark.Bool {
	if f.Fs == nil || f.Path == "" {
		return starlark.False
	}
	return starlark.True
}

// ParseParam from starlark input
func ParseParam(parser script.ArgParser, index int) (*File, error) {
	val, err := parser.GetParam(index)
	if err != nil {
		return nil, err
	}

	if val == nil {
		return nil, fmt.Errorf("expected file type, got None")
	}

	f, ok := val.(*File)
	if !ok {
		return nil, fmt.Errorf("%w: expected file type", script.ErrInvalidArgType)
	}

	return f, nil
}
