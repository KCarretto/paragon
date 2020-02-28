package regex

import (
	"github.com/kcarretto/paragon/pkg/script"
)

// Library prepares a new regex library for use within a script environment.
func Library() script.Library {
	return script.Library{
		"replace": script.Func(replace),
	}
}

// Include the sys library in a script environment.
func Include() script.Option {
	return script.WithLibrary("regex", Library())
}
