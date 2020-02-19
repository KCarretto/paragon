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
	Sync() error
}

// Library prepares a new file library for use within a script environment.
func Library() script.Library {
	return script.Library{
		"move":    script.Func(move),
		"name":    script.Func(name),
		"content": script.Func(content),
		"write":   script.Func(write),
		"copy":    script.Func(copy),
		"remove":  script.Func(remove),
		"chown":   script.Func(chown),
		"chmod":   script.Func(chmod),
	}
}

// Include the file library in a script environment.
func Include() script.Option {
	return script.WithLibrary("file", Library())
}
