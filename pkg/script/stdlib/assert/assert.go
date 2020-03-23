package assert

import (
	"github.com/kcarretto/paragon/pkg/script"
)

// Library prepares a new regex library for use within a script environment.
func Library() script.Library {
	return script.Library{
		"noError": script.Func(noError),
		"equal":   script.Func(equal),
	}
}

// Include the sys library in a script environment.
func Include() script.Option {
	return script.WithLibrary("assert", Library())
}
