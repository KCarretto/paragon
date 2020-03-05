// +build dev

package drop

// deleteFile is a no-op for dev mode
func deleteFile(path string) error {
	path = filepath.Clean(path)
	log.Printf("[INFO] would have attempted deletion of file: %q", path)
	return nil
}
