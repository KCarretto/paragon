package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/kcarretto/paragon/pkg/script"
	libassert "github.com/kcarretto/paragon/pkg/script/stdlib/assert"
	libassets "github.com/kcarretto/paragon/pkg/script/stdlib/assets"
	libcrypto "github.com/kcarretto/paragon/pkg/script/stdlib/crypto"
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

func run(ctx context.Context, assets http.FileSystem) error {
	env := &libenv.Environment{}

	assetEnv := &libassets.Environment{
		Assets: assets,
	}
	task, err := assetEnv.OpenFile("scripts/main.rg")
	if err != nil {
		return fmt.Errorf("failed to execute scripts/main.rg: %w", err)
	}

	code := script.New("main.rg", task, script.WithOutput(os.Stdout),
		env.Include(),
		assetEnv.Include(),
		libassert.Include(),
		libcrypto.Include(),
		libfile.Include(),
		libhttp.Include(),
		libproc.Include(),
		libregex.Include(),
		libsys.Include(),
	)

	return code.Exec(ctx)
}

func main() {
	var bundlePath string
	var bundleKey string

	app := &cli.App{
		Name:  "renegade",
		Usage: "Interpreter for the renegade scripting language.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "bundle",
				Usage:       "Path to the bundle tar.gz file.",
				Destination: &bundlePath,
			},
			&cli.StringFlag{
				Name:        "key",
				Usage:       "Base64 representation of Key to be used to decrypt bundle.",
				Destination: &bundleKey,
			},
		},
		Action: func(c *cli.Context) error {
			var bundle http.FileSystem
			if bundlePath != "" {
				bundleFile, err := ioutil.ReadFile(bundlePath)
				if err != nil {
					return fmt.Errorf("failed to open bundle file %q: %w", bundlePath, err)
				}

				if bundleKey != "" {
					cryptoKey, err := libcrypto.CreateKey(bundleKey)
					if err != nil {
						return fmt.Errorf("failed to import key: %w", err)
					}
					decryptedBundle, err := libcrypto.Decrypt(cryptoKey, string(bundleFile))
					if err != nil {
						return fmt.Errorf("failed to decrypt bundle: %w", err)
					}
					bundleFile = []byte(decryptedBundle)
				}

				assetTar := &libassets.TarGZBundler{
					Buffer: bytes.NewBuffer(bundleFile),
				}
				bundle, err = assetTar.FileSystem()
				if err != nil {
					return fmt.Errorf("failed to create assets filesystem: %w", err)
				}
			}

			return run(context.Background(), bundle)

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
