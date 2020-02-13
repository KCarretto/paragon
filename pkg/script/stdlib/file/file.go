// Package file provides the standard file library used to perform common file operations.
package file

import (
	"io"
	"os"

	"github.com/kcarretto/paragon/pkg/script"
)

// File defines the required methods for a file object to be used for file operations.
type File interface {
	io.Reader
	io.Writer

	Name() string
	Chmod(os.FileMode) error
	Chown(uid, gid int) error
	Stat() (os.FileInfo, error)

	Move(dstPath string) error
	Remove() error
}

// file wraps a provided File implementation to implement a starlark.Value
// type file struct {
// 	f File
// }

// Import the sys library to enable scripts to access low level system functionality.
func Import() script.Library {
	return script.Library{
		// openFile
		// exists
		// name
		// path
		// content
		// write
		// replace
		// copy
		// move
		// remove
		// setOwner
		// setWrite
		// setRead
		// setExec
		// setSUID
	}
}
