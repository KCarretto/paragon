package assets

import (
	"fmt"
	"net/http"
	"os"
)

var ErrUnsupported = fmt.Errorf("operation not supported for asset file")

type File struct {
	http.File

	name string
}

// Name of the file as provided to openFile.
func (f File) Name() string {
	return f.name
}

// Write is unsupported and always returns an error.
func (f File) Write(p []byte) (int, error) {
	return 0, ErrUnsupported
}

// Chmod is unsupported and always returns an error.
func (f File) Chmod(os.FileMode) error {
	return ErrUnsupported
}

// Chown is unsupported and always returns an error.
func (f File) Chown(uid, gid int) error {
	return ErrUnsupported
}

// Move is unsupported and always returns an error.
func (f File) Move(dstPath string) error {
	return ErrUnsupported
}

// Remove is unsupported and always returns an error.
func (f File) Remove() error {
	return ErrUnsupported
}
