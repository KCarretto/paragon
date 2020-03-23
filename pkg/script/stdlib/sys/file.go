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
func (f *File) Move(dstPath string) error {
	dir := path.Dir(dstPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create parent directory %q: %w", dir, err)
	}

	// Must first close file, otherwise on windows the following error occurs:
	// The process cannot access the file because it is being used by another process.
	f.File.Close()

	err := os.Rename(f.Name(), dstPath)
	if err != nil {
		return err
	}

	// Re-open destination file after moving
	newF, err := os.Open(dstPath)
	if err != nil {
		return err
	}

	f.File = newF
	return nil
}

// Remove uses os.Remove to remove a file/folder. WARNING: basically works like rm -rf
func (f *File) Remove() error {
	f.Close()
	return os.RemoveAll(f.Name())
}

// Sync closes and then reopens the file.
func (f *File) Sync() error {
	if f == nil || f.File == nil {
		return nil
	}
	f.File.Close()

	newF, err := os.Open(f.Name())
	if err != nil {
		return err
	}

	f.File = newF
	return nil
}
