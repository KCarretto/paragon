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
