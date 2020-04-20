// +build !dev

package drop

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// DeleteFile is used to delete the running executable with the provided path.
func DeleteFile(path string) error {
	path = filepath.Clean(path)

	log.Printf("[INFO] attempting deletion of file: %q", path)
	if path == "" || path == string(filepath.Separator) {
		return fmt.Errorf("file resolved to an invalid path (empty or root directory)")
	}

	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("failed to remove file: %w", err)
	}

	return nil
}
