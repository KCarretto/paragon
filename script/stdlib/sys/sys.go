package sys

import (
	"github.com/kcarretto/paragon/script"
)

// Lib is a map that has keys as symbols to be used in a starlark script and values to functions that will be run.
var Lib script.Library = script.Library{
	"move":          script.Func(Move),
	"copy":          script.Func(Copy),
	"remove":        script.Func(Remove),
	"exec":          script.Func(Exec),
	"read":          script.Func(ReadFile),
	"write":         script.Func(WriteFile),
	"chmod":         script.Func(Chmod),
	"chown":         script.Func(Chown),
	"processes":     script.Func(Processes),
	"kill":          script.Func(Kill),
	"connections":   script.Func(Connections),
	"dir":           script.Func(Dir),
	"replaceString": script.Func(ReplaceString),
	"request":       script.Func(Request),
	"detectOS":      script.Func(Detect),
}
