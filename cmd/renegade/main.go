package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/kcarretto/paragon/pkg/script"
	libassert "github.com/kcarretto/paragon/pkg/script/stdlib/assert"
	libassets "github.com/kcarretto/paragon/pkg/script/stdlib/assets"
	libenv "github.com/kcarretto/paragon/pkg/script/stdlib/env"
	libfile "github.com/kcarretto/paragon/pkg/script/stdlib/file"
	libhttp "github.com/kcarretto/paragon/pkg/script/stdlib/http"
	libproc "github.com/kcarretto/paragon/pkg/script/stdlib/process"
	libregex "github.com/kcarretto/paragon/pkg/script/stdlib/regex"
	libsys "github.com/kcarretto/paragon/pkg/script/stdlib/sys"

	"github.com/urfave/cli"
	"go.starlark.net/starlark"
)

func compilePredeclared(libs map[string]script.Library) starlark.StringDict {
	builtins := make(starlark.StringDict)

	for name, lib := range libs {
		builtins[name] = lib
	}

	return builtins
}

func run(ctx context.Context, id string, content io.Reader, assets http.FileSystem) error {
	env := &libenv.Environment{}

	assetEnv := &libassets.Environment{
		Assets: assets,
	}

	code := script.New(id, content, script.WithOutput(os.Stdout),
		env.Include(),
		assetEnv.Include(),
		libassert.Include(),
		libfile.Include(),
		libhttp.Include(),
		libproc.Include(),
		libregex.Include(),
		libsys.Include(),
	)

	return code.Exec(ctx)
}

func main() {
	var (
		assetPath string
		taskPath  string
		id        string
		// key       string
	)

	app := &cli.App{
		Name:  "renegade",
		Usage: "Interpreter for the renegade scripting language.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "id",
				Usage:       "ID to associate with the task.",
				Destination: &id,
			},
			&cli.StringFlag{
				Name:        "assets",
				Usage:       "Path to the assets tar file.",
				Destination: &assetPath,
			},
			&cli.StringFlag{
				Name:        "task",
				Usage:       "Path to the task file.",
				Destination: &taskPath,
			},
		},
		Action: func(c *cli.Context) error {
			var assets http.FileSystem
			if assetPath != "" {
				assetFile, err := ioutil.ReadFile(assetPath)
				if err != nil {
					return fmt.Errorf("failed to open assets file %q: %w", assetPath, err)
				}
				assetTar := &libassets.TarGZBundler{
					Buffer: bytes.NewBuffer(assetFile),
				}
				assets, err = assetTar.FileSystem()
				if err != nil {
					return fmt.Errorf("failed to create assets filesystem: %w", err)
				}
			}

			var task io.Reader = os.Stdin
			if taskPath != "" {
				taskFile, err := ioutil.ReadFile(taskPath)
				if err != nil {
					return fmt.Errorf("failed to open task file %q: %w", taskPath, err)
				}

				task = bytes.NewBuffer(taskFile)
			}

			return run(context.Background(), id, task, assets)

		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	// thread := &starlark.Thread{Load: repl.MakeLoad()}
	// builtins := compilePredeclared(
	// 	map[string]script.Library{
	// 		"sys":     sys.Library(),
	// 		"http":    http.Library(),
	// 		"file":    file.Library(),
	// 		"regex":   regex.Library(),
	// 		"process": process.Library(),
	// 		"assert":  assert.Library(),
	// 	},
	// )
	// switch len(os.Args) {
	// case 1:
	// 	fmt.Println("Welcome to Renegade Shell!")
	// 	repl.REPL(thread, builtins)
	// case 2:
	// 	// Execute specified file.
	// 	filename := os.Args[1]
	// 	var err error
	// 	_, err = starlark.ExecFile(thread, filename, nil, builtins)
	// 	if err != nil {
	// 		repl.PrintError(err)
	// 		os.Exit(1)
	// 	}
	// default:
	// 	log.Fatal("want at most one renegade file name")
	// }

}
