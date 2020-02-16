package sys

import (
	"fmt"
	"os"
	"path"
)

// File is our impl of File.
type File struct {
	*os.File
}

// Move uses os.Rename to move a file/folder
func (f File) Move(dstPath string) error {
	dir := path.Dir(dstPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create parent directory %q: %w", dir, err)
	}

	err := os.Rename(f.Name(), dstPath)
	if err != nil {
		return err
	}

	err = f.File.Close()
	if err != nil {
		return err
	}

	newF, err := os.Open(dstPath)
	if err != nil {
		return err
	}

	f.File = newF
	return nil
}

// Remove uses os.Remove to remove a file/folder. WARNING: basically works like rm -rf
func (f File) Remove() error {
	f.Close()
	return os.RemoveAll(f.Name())
}
