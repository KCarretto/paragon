package sys

import (
	"github.com/kcarretto/paragon/pkg/script"
)

// Library prepares a new sys library for use within a script environment.
func Library() script.Library {
	return script.Library{
		"openFile":    script.Func(openFile),
		"detectOS":    script.Func(detectOS),
		"exec":        script.Func(exec),
		"connections": script.Func(connections),
		"processes":   script.Func(processes),
	}
}

// @DEPRECATE Import the sys library to enable scripts to access low level system functionality.
func Import() script.Library {
	return Library()
}

// Include the sys library in a script environment.
func Include() script.Option {
	return script.WithLibrary("sys", Library())
}
