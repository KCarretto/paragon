package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
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

	"github.com/spf13/afero"
	"github.com/urfave/cli"
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

func run(ctx context.Context, assets afero.Fs) error {
	env := &libenv.Environment{}

	assetEnv := &libassets.Environment{
		Assets: assets,
	}
	task, err := assetEnv.Assets.Open("scripts/main.rg")
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
	passedArgs := len(os.Args) != 1

	if passedArgs {
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
				var bundle afero.Fs
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
	} else {
		fmt.Println("Welcome to Renegade Shell!")
		thread := &starlark.Thread{Load: repl.MakeLoad()}
		builtins := compilePredeclared(
			map[string]script.Library{
				"sys":     libsys.Library(),
				"http":    libhttp.Library(),
				"file":    libfile.Library(),
				"regex":   libregex.Library(),
				"process": libproc.Library(),
				"assert":  libassert.Library(),
			},
		)
		repl.REPL(thread, builtins)
	}
}
