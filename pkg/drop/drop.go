// Package drop provides a simple method to walk a filesystem and execute scripts that will drop...
package drop

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/kcarretto/paragon/pkg/script"
	"github.com/kcarretto/paragon/pkg/script/stdlib"

	"github.com/shurcooL/httpfs/vfsutil"
)

// TheBase uses WubWubWubWUBWUBWUBWUB.
func TheBase(ctx context.Context, assets http.FileSystem) {
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

		// Initialize script
		dropScript := script.New(fi.Name(), content,
			script.WithOutput(os.Stdout), // TODO: Output to stdout?
			stdlib.Load(stdlib.WithAssets(assets)),
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
