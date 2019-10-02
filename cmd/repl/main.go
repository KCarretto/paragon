package main

import (
	"fmt"

	"github.com/kcarretto/paragon/script"
	"github.com/kcarretto/paragon/script/stdlib"
	"go.starlark.net/repl"
	"go.starlark.net/starlark"
)

func main() {
	script.NewInterpreter()

	thread := starlark.Thread{
		Print: func(_ *starlark.Thread, msg string) { fmt.Println(msg) },
		Load:  stdlib.Loader(),
	}

	repl.REPL(&thread, starlark.StringDict{})
}
