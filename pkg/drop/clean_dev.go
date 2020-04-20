// +build dev

package drop

import (
	"log"
	"path/filepath"
)

// DeleteFile is a no-op for dev mode
func DeleteFile(path string) error {
	path = filepath.Clean(path)
	log.Printf("[INFO] would have attempted deletion of file: %q", path)
	return nil
}
