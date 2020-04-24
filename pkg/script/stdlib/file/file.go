// Package file provides the standard file library used to perform common file operations.
package file

import (
	"github.com/kcarretto/paragon/pkg/script"
)

// Library prepares a new file library for use within a script environment.
func Library() script.Library {
	return script.Library{
		"move":    script.Func(move),
		"name":    script.Func(name),
		"content": script.Func(content),
		"write":   script.Func(write),
		"copy":    script.Func(copy),
		"remove":  script.Func(remove),
		"chmod":   script.Func(chmod),
		"hash":    script.Func(hash),
		"exists":  script.Func(exists),
		"drop":    script.Func(drop),
	}
}

// Include the file library in a script environment.
func Include() script.Option {
	return script.WithLibrary("file", Library())
}
