package http

import "github.com/kcarretto/paragon/pkg/script"

// Library prepares a new http library for use within a script environment.
func Library() script.Library {
	return script.Library{
		"newRequest": script.Func(newRequest),
		"setMethod":  script.Func(setMethod),
		"setHeader":  script.Func(setHeader),
		"setBody":    script.Func(setBody),
		"exec":       script.Func(exec),
	}
}

// Include the http library in a script environment.
func Include() script.Option {
	return script.WithLibrary("http", Library())
}
