package enum

import (
	"github.com/kcarretto/paragon/pkg/script"
)

// Library prepares a new enumeration library for use within a script environment.
func Library() script.Library {
	return script.Library{
		"scan": script.Func(scan),
	}
}

// Include the sys library in a script environment.
func Include() script.Option {
	return script.WithLibrary("enum", Library())
}
