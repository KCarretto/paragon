package asset

import (
	"github.com/kcarretto/paragon/script"
)

// Lib is a map that has keys as symbols to be used in a starlark script and values to functions that will be run.
var Lib script.Library = script.Library{
	"load": script.Func(Load),
}
