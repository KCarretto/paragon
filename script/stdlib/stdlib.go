// Package stdlib contains all of the basic building blocks of OS operation. It is intended to be cross platform and
// minimalistic in scale. Its largest sublibrary is "sys" and to use it simply use the syntax `load("sys", "myFunc")`
// where `myFunc` is one of the premade library functions.
package stdlib

import (
	"github.com/kcarretto/paragon/script"
	"github.com/kcarretto/paragon/script/stdlib/sys"
)

var libs = map[string]script.Library{
	"sys": sys.Lib,
}

// Load is a script Option that loads the standard library into a script's execution environment
func Load() script.Option {
	return script.WithLibraries(libs)
}
