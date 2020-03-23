package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kcarretto/paragon/pkg/script"
	"github.com/kcarretto/paragon/pkg/script/stdlib/assert"
	"github.com/kcarretto/paragon/pkg/script/stdlib/file"
	"github.com/kcarretto/paragon/pkg/script/stdlib/http"
	"github.com/kcarretto/paragon/pkg/script/stdlib/process"
	"github.com/kcarretto/paragon/pkg/script/stdlib/regex"
	"github.com/kcarretto/paragon/pkg/script/stdlib/sys"
	"go.starlark.net/repl"
	"go.starlark.net/starlark"
)

func compilePredeclared(libs map[string]script.Library) starlark.StringDict {
	builtins := make(starlark.StringDict)

	for name, lib := range libs {
		builtins[name] = lib
	}

	return builtins
}

func main() {
	thread := &starlark.Thread{Load: repl.MakeLoad()}
	builtins := compilePredeclared(
		map[string]script.Library{
			"sys":     sys.Library(),
			"http":    http.Library(),
			"file":    file.Library(),
			"regex":   regex.Library(),
			"process": process.Library(),
			"assert":  assert.Library(),
		},
	)
	switch len(os.Args) {
	case 1:
		fmt.Println("Welcome to Renegade Shell!")
		repl.REPL(thread, builtins)
	case 2:
		// Execute specified file.
		filename := os.Args[1]
		var err error
		_, err = starlark.ExecFile(thread, filename, nil, builtins)
		if err != nil {
			repl.PrintError(err)
			os.Exit(1)
		}
	default:
		log.Fatal("want at most one renegade file name")
	}

}
