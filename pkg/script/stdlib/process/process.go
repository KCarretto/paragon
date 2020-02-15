package process

import "github.com/kcarretto/paragon/pkg/script"

// Library prepares a new process library for use within a script environment.
func Library() script.Library {
	return script.Library{
		"kill": script.Func(kill),
	}
}

// Include the process library in a script environment.
func Include() script.Option {
	return script.WithLibrary("process", Library())
}
