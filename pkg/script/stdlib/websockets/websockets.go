package websockets

import (
	"github.com/kcarretto/paragon/pkg/script"
)

// Library prepares a new enumeration library for use within a script environment.
func Library() script.Library {
	return script.Library{
		"giveshell": script.Func(giveshell),
	}
}

// Include the sys library in a script environment.
func Include() script.Option {
	return script.WithLibrary("websockets", Library())
}
