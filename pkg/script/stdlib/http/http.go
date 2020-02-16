package http

import "github.com/kcarretto/paragon/pkg/script"

// Library prepares a new process library for use within a script environment.
func Library() script.Library {
	return script.Library{
		"newRequest": script.Func(newRequest),
		"setMethod":  script.Func(setMethod),
		"setHeader":  script.Func(setHeader),
		"setBody":    script.Func(setBody),
		"exec":       script.Func(exec),
	}
}

// Include the process library in a script environment.
func Include() script.Option {
	return script.WithLibrary("process", Library())
}
