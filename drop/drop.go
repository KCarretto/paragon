// Package drop provides a simple method to walk a filesystem and execute scripts that will drop...
package drop

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/kcarretto/paragon/script"
	"github.com/kcarretto/paragon/script/stdlib"
	"google.golang.org/appengine/log"

	"github.com/shurcooL/httpfs/vfsutil"
)

// TheBase uses WubWubWubWUBWUBWUBWUB.
func TheBase(ctx context.Context, assets http.FileSystem) {

	if err := vfsutil.Walk(assets, "/scripts", func(path string, fi os.FileInfo, err error) error {
		// Check for stat error
		if err != nil {
			log.Errorf(ctx, "failed to stat file %q: %s", path, err.Error())
			return nil
		}

		// Skip directories
		if fi.IsDir() {
			return nil
		}

		// Skip files that don't end with .rg
		if !strings.HasSuffix(fi.Name(), ".rg") {
			log.Warningf(context.Background(), "found non-script in scripts directory, scripts must end in '.rg' %q: %s", path, err.Error())
			return nil
		}

		// Open script file
		content, err := assets.Open(path)
		if err != nil {
			log.Errorf(ctx, "failed to open script %q: %s", path, err.Error())
			return nil
		}

		// Initialize script
		dropScript := script.New(fi.Name(), content,
			script.WithOutput(os.Stdout), // TODO: Output to stdout?
			stdlib.Load(stdlib.WithAssets(assets)),
		)

		// Run script
		if err := dropScript.Exec(context.Background()); err != nil {
			log.Errorf(ctx, "script failed execution %q: %s", path, err.Error())
			return nil
		}

		return nil
	}); err != nil {
		log.Errorf(ctx, "failed to walk files: %s", err.Error())
	}
}
